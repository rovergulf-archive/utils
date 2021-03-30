package pgxs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"io/ioutil"
)

func newConn(ctx context.Context, lg *zap.SugaredLogger, conf *Config) (*pgx.Conn, error) {
	if conf == nil {
		return nil, fmt.Errorf("pxgs: Config should'n be nil")
	}

	if conf.TLS.Enabled {
		lg.Debugf("Client TLS connection enabled")
		conf.TLSConfig = new(tls.Config)
		if len(conf.TLS.CaPath) > 0 {
			caCert, err := ioutil.ReadFile(conf.TLS.CaPath)
			if err != nil {
				return nil, fmt.Errorf("pgxs: Unable to load CA cert: %s", err)
			}
			caCertPool := x509.NewCertPool()

			caCertPool.AppendCertsFromPEM(caCert)

			conf.TLSConfig.ClientCAs = caCertPool
			conf.TLSConfig.InsecureSkipVerify = conf.TLS.Verify
		}

		if len(conf.TLS.CertPath) > 0 && len(conf.TLS.KeyPath) > 0 {
			cert, err := tls.LoadX509KeyPair(conf.TLS.CertPath, conf.TLS.KeyPath)
			if err != nil {
				return nil, fmt.Errorf("pgxs: Unable to load tls keypair: %s", err)
			}
			conf.TLSConfig.Certificates = append(conf.TLSConfig.Certificates, cert)
		}
	} else {
		lg.Debugf("Client TLS connection disabled")
	}

	conn, err := pgx.Connect(ctx, conf.GetConnString())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func connFromString(ctx context.Context, connString string) (*pgconn.PgConn, error) {
	return pgconn.Connect(ctx, connString)
}

func (db *Repo) ConnectDB(ctx context.Context, tlsConfig *tls.Config) (*pgconn.PgConn, error) {
	conf, err := db.GetConnConfig()
	if err != nil {
		return nil, fmt.Errorf("pgxs: Unable to prepare postgres config: %s", err)
	}
	conf.TLSConfig = tlsConfig

	return pgconn.ConnectConfig(ctx, conf)
}
