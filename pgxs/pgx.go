package pgxs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	ConnString      string      `json:"conn_string" yaml:"conn_string"`
	ActualSchema    string      `json:"actual_schema" yaml:"schema_path"`
	MigrationSchema string      `json:"migration_schema" yaml:"migration_schema_path"`
	DataDir         string      `json:"data_dir" yaml:"data_dir"`
	Host            string      `json:"host" yaml:"host"`
	Port            string      `json:"port" yaml:"port"`
	Name            string      `json:"name" yaml:"name"`
	User            string      `json:"user" yaml:"user"`
	Password        string      `json:"password" yaml:"password"`
	SslMode         string      `json:"ssl_mode" yaml:"ssl_mode"`
	TLS             SSL         `json:"tls" yaml:"tls"`
	TLSConfig       *tls.Config `json:"-" yaml:"-"`
}

type SSL struct {
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	Verify   bool   `json:"verify" yaml:"verify"`
	CaPath   string `json:"ca" yaml:"ca"`
	KeyPath  string `json:"key" yaml:"key"`
	CertPath string `json:"cert" yaml:"cert"`
}

type Repo struct {
	Logger *zap.SugaredLogger `json:"-" yaml:"-"`
	Pool   *pgxpool.Pool      `json:"-" yaml:"-"`
	Config *Config            `json:"-" yaml:"-"`
}

func New(ctx context.Context, lg *zap.SugaredLogger, conf *Config) (*Repo, error) {
	if conf == nil {
		return nil, fmt.Errorf("pxgs: Config should'n be nil")
	}

	s := &Repo{
		Logger: lg.Named("pgxs"),
		Config: conf,
	}

	if s.Config.TLS.Enabled {
		s.Logger.Debugf("Client TLS connection enabled")
		s.Config.TLSConfig = new(tls.Config)
		if len(s.Config.TLS.CaPath) > 0 {
			caCert, err := ioutil.ReadFile(s.Config.TLS.CaPath)
			if err != nil {
				return nil, fmt.Errorf("pgxs: Unable to load CA cert: %s", err)
			}
			caCertPool := x509.NewCertPool()

			caCertPool.AppendCertsFromPEM(caCert)

			s.Config.TLSConfig.ClientCAs = caCertPool
			s.Config.TLSConfig.InsecureSkipVerify = s.Config.TLS.Verify
		}

		if len(s.Config.TLS.CertPath) > 0 && len(s.Config.TLS.KeyPath) > 0 {
			cert, err := tls.LoadX509KeyPair(s.Config.TLS.CertPath, s.Config.TLS.KeyPath)
			if err != nil {
				return nil, fmt.Errorf("pgxs: Unable to load tls keypair: %s", err)
			}
			s.Config.TLSConfig.Certificates = append(s.Config.TLSConfig.Certificates, cert)
		}
	} else {
		s.Logger.Debugf("Client TLS connection disabled")
	}

	pool, err := s.ConnectDB(ctx, s.Config.TLSConfig)
	if err != nil {
		return nil, err
	}

	s.Pool = pool

	return s, nil
}

func (db *Repo) PoolFromString(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, connString)
}

//
func (db *Repo) ConnectDB(ctx context.Context, tlsConfig *tls.Config) (*pgxpool.Pool, error) {
	conf, err := db.GetPgxConfig()
	if err != nil {
		return nil, fmt.Errorf("pgxs: Unable to prepare postgres config: %s", err)
	}
	conf.ConnConfig.TLSConfig = tlsConfig

	return pgxpool.ConnectConfig(ctx, conf)
}

func (db *Repo) GetPgxConnString() string {
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

func (db *Repo) GetPgxConfig() (*pgxpool.Config, error) {
	c, err := pgxpool.ParseConfig(db.GetPgxConnString())
	if err != nil {
		return nil, fmt.Errorf("pgxs: unable to parse pgx config: %s", err)
	}

	return c, nil
}

func (db *Repo) GracefulShutdown() {
	if db.Pool != nil {
		db.Pool.Close()
		db.Logger.Infof("Successfully closed postgreSQL connection pool")
	}
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

// handleSqlErr used to avoid not exists and already exists debug queries
func (db *Repo) DebugLogSqlErr(q string, err error) error {
	pgErr, deuce := err.(*pgconn.PgError)
	if deuce {
		if pgErr.Code == "23505" {
			deuce = false
		}
	}

	if err != pgx.ErrNoRows && !deuce {
		db.Logger.Debugf("query: \n%s", q)
	}

	return err
}
