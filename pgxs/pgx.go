package pgxs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type PgConf struct {
	ConnString       string      `json:"conn_string" yaml:"conn_string"`
	ServiceName      string      `json:"service" yaml:"service"`
	ActualSchemaPath string      `json:"actual_schema_path" yaml:"schema_path"`
	DataDir          string      `json:"data_dir" yaml:"data_dir"`
	Host             string      `json:"host" yaml:"host"`
	Port             string      `json:"port" yaml:"port"`
	Name             string      `json:"name" yaml:"name"`
	User             string      `json:"user" yaml:"user"`
	Password         string      `json:"password" yaml:"password"`
	SslMode          string      `json:"ssl_mode" yaml:"ssl_mode"`
	SslPath          string      `json:"ssl_path" yaml:"ssl_path"`
	TLS              SSL         `json:"tls" yaml:"tls"`
	TLSConfig        *tls.Config `json:"-" yaml:"-"`
}

type SSL struct {
	CaPath   string `json:"ca" yaml:"ca"`
	KeyPath  string `json:"key" yaml:"key"`
	CertPath string `json:"cert" yaml:"cert"`
	Insecure bool   `json:"insecure" yaml:"insecure"`
}

type Repo struct {
	ServiceName string             `json:"service" yaml:"service"`
	Logger      *zap.SugaredLogger `json:"-" yaml:"-"`
	Pool        *pgxpool.Pool      `json:"-" yaml:"-"`
	Config      *PgConf
	connected   bool
}

func (db *Repo) PoolFromString(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, connString)
}

//
func (db *Repo) ConnectDB(ctx context.Context, tlsConfig *tls.Config) (*pgxpool.Pool, error) {
	conf, err := db.GetConfig()
	if err != nil {
		db.Logger.Errorf("Unable to prepare postgres config: %s", err)
		return nil, err
	}
	conf.ConnConfig.TLSConfig = tlsConfig

	return pgxpool.ConnectConfig(ctx, conf)
}

func (db *Repo) GetPgxPoolConnString() string {
	dbUrl := os.Getenv("DB_URL")
	if len(dbUrl) > 0 {
		return dbUrl
	}

	connString := fmt.Sprintf("host=%s port=%s database=%s user=%s password=%s sslmode=%s",
		db.Config.Host,
		db.Config.Port,
		db.Config.Name,
		db.Config.User,
		db.Config.Password,
		db.Config.SslMode,
	)

	if len(db.Config.TLS.CaPath) > 0 {
		connString += fmt.Sprintf(" sslrootcert=%s", db.Config.TLS.CaPath)
	}

	if len(db.Config.TLS.CertPath) > 0 {
		connString += fmt.Sprintf(" sslcert=%s", db.Config.TLS.CertPath)
	}

	if len(db.Config.TLS.KeyPath) > 0 {
		connString += fmt.Sprintf(" sslkey=%s", db.Config.TLS.KeyPath)
	}

	return connString
}

func (db *Repo) GetPgxPoolString() string {
	return fmt.Sprintf("host=%s port=%s database=%s user=%s password=%s sslmode=%s",
		db.Config.Host,
		db.Config.Port,
		db.Config.Name,
		db.Config.User,
		db.Config.Password,
		db.Config.SslMode)
}

func (db *Repo) GetPgxConfig() (*pgxpool.Config, error) {
	c, err := pgxpool.ParseConfig(db.GetPgxPoolString())
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (db *Repo) GetConfig() (*pgxpool.Config, error) {
	c, err := db.GetPgxConfig()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (db *Repo) GracefulShutdown(ctx context.Context) {
	if db.Pool != nil {
		db.Pool.Close()
		log.Println("Successfully closed postgreSQL connection pool")
	}
}

func NewSecurePool(ctx context.Context, lg *zap.SugaredLogger, serviceName string, pgconf *PgConf) (Repo, error) {
	if pgconf == nil {
		return Repo{}, fmt.Errorf("%s", "Config should'n be nil")
	}

	s := Repo{
		Logger:      lg,
		ServiceName: serviceName,
		Config:      pgconf,
	}

	conf := new(tls.Config)

	caFile := path.Join(s.Config.SslPath, s.Config.TLS.CaPath)
	if len(caFile) > len(s.Config.SslPath) {
		s.Logger.Debugf("Client CA File would be used for TLS connection")
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			s.Logger.Fatal(err)
		}
		caCertPool := x509.NewCertPool()

		caCertPool.AppendCertsFromPEM(caCert)

		conf.ClientCAs = caCertPool
		conf.InsecureSkipVerify = s.Config.TLS.Insecure
	}

	pool, err := s.ConnectDB(ctx, conf)
	if err != nil {
		s.Logger.Errorf("Unable to connect psql intsance: %s", err)
		return s, err
	}

	s.Pool = pool

	return s, nil
}

func NewPool(ctx context.Context, lg *zap.SugaredLogger, serviceName string, pgconf *PgConf) (Repo, error) {
	if pgconf == nil {
		return Repo{}, fmt.Errorf("%s", "Config should'n be nil")
	}

	s := Repo{
		Logger:      lg,
		ServiceName: serviceName,
		Config:      pgconf,
	}

	pool, err := s.ConnectDB(ctx, nil)
	if err != nil {
		s.Logger.Errorf("Unable to connect psql intsance: %s", err)
		return s, err
	}
	s.Pool = pool

	return s, nil
}

//
// according to https://github.com/jackc/pgx/blob/master/conn.go#L84
// have to watch changes, to prevent internal issues
//
func QuoteString(str string) string {
	str = strings.Replace(str, "'", "", -1)
	str = strings.Replace(str, "%", "", -1)
	return str
}

func (db *Repo) SanitizeString(str string) string {
	return QuoteString(str)
}
