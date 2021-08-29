package pgxs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"io/ioutil"
)

func NewPool(ctx context.Context, lg *zap.SugaredLogger, conf *Config) (*Repo, error) {
	if conf == nil {
		return nil, ErrEmptyConfig
	}

	s := &Repo{
		Logger: lg.Named("pgx_pool"),
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

	pool, err := s.ConnectDBPool(ctx, s.Config.TLSConfig)
	if err != nil {
		return nil, err
	}

	s.Pool = pool

	return s, nil
}

func (db *Repo) PoolFromString(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, connString)
}

// ConnectDBPool initializes Pool connection
func (db *Repo) ConnectDBPool(ctx context.Context, tlsConfig *tls.Config) (*pgxpool.Pool, error) {
	conf, err := db.GetPoolConfig()
	if err != nil {
		return nil, fmt.Errorf("pgxs: Unable to prepare postgres config: %s", err)
	}
	conf.ConnConfig.TLSConfig = tlsConfig

	return pgxpool.ConnectConfig(ctx, conf)
}
