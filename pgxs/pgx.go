package pgxs

import (
	"crypto/tls"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	DbUri          string      `json:"db_uri" yaml:"db_uri"`
	MigrationsPath string      `json:"migration_schemas" yaml:"migration_schemas"`
	DataDir        string      `json:"data_dir" yaml:"data_dir"`
	Host           string      `json:"host" yaml:"host"`
	Port           string      `json:"port" yaml:"port"`
	Name           string      `json:"name" yaml:"name"`
	User           string      `json:"user" yaml:"user"`
	Password       string      `json:"password" yaml:"password"`
	SslMode        string      `json:"ssl_mode" yaml:"ssl_mode"`
	TLS            SSL         `json:"tls" yaml:"tls"`
	TLSConfig      *tls.Config `json:"-" yaml:"-"`
}

type SSL struct {
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	Verify   bool   `json:"verify" yaml:"verify"`
	CaPath   string `json:"ca" yaml:"ca"`
	KeyPath  string `json:"key" yaml:"key"`
	CertPath string `json:"cert" yaml:"cert"`
}

func (c *Config) GetConnString() string {
	dbUrl := os.Getenv("DB_URL")
	if len(dbUrl) > 0 {
		return dbUrl
	}

	connString := fmt.Sprintf("host=%s port=%s database=%s user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Name,
		c.User,
		c.Password,
		c.SslMode,
	)

	if len(c.TLS.CaPath) > 0 {
		connString += fmt.Sprintf(" sslrootcert=%s", c.TLS.CaPath)
	}

	if len(c.TLS.CertPath) > 0 {
		connString += fmt.Sprintf(" sslcert=%s", c.TLS.CertPath)
	}

	if len(c.TLS.KeyPath) > 0 {
		connString += fmt.Sprintf(" sslkey=%s", c.TLS.KeyPath)
	}

	return connString
}

type Repo struct {
	Logger *zap.SugaredLogger `json:"-" yaml:"-"`
	Pool   *pgxpool.Pool      `json:"-" yaml:"-"`
	Config *Config            `json:"-" yaml:"-"`
}

func (db *Repo) GetConnConfig() (*pgconn.Config, error) {
	c, err := pgconn.ParseConfig(db.Config.GetConnString())
	if err != nil {
		return nil, fmt.Errorf("pgxs: unable to parse pgx config: %s", err)
	}

	return c, nil
}

func (db *Repo) GetPoolConfig() (*pgxpool.Config, error) {
	c, err := pgxpool.ParseConfig(db.Config.GetConnString())
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
