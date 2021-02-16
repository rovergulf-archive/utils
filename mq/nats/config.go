package natsmq

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"io/ioutil"
)

// Config represent stan.Conn and nats.Conn connection parameters
type Config struct {
	ClientId  string             `json:"client_id" yaml:"client_id"`
	ClusterId string             `json:"clusterId" yaml:"cluster_id"`
	BrokerId  string             `json:"brokerId" yaml:"broker_id"`
	Broker    string             `json:"brokers" yaml:"broker"`
	User      string             `json:"user" yaml:"user"`
	Password  string             `json:"password" yaml:"password"`
	Token     string             `json:"token" yaml:"token"`
	Channels  map[string]string  `json:"channels" yaml:"channels"`
	TLS       SSL                `json:"tls" yaml:"tls"`
	NatsConn  []nats.Option      `json:"-" yaml:"-"`
	StanConn  []stan.Option      `json:"-" yaml:"-"`
	Logger    *zap.SugaredLogger `json:"-" yaml:"-"`
}

func (c *Config) GetNatsUserInfo() nats.Option {
	return nats.UserInfo(c.User, c.Password)
}

type NatsSubOpts struct {
	Sub     nats.SubscriptionType `json:"-" yaml:"-"`
	Subject string                `json:"subject" yaml:"subject"`
	*Config
}

type StanSubOpts struct {
	Opts    []stan.SubscriptionOption `json:"-" yaml:"-"`
	Channel string                    `json:"channel" yaml:"channel"`
	*Config
}

// SSL contains tls connection options
type SSL struct {
	Enabled  bool               `json:"enabled" yaml:"enabled"`
	CaPath   string             `json:"ca" yaml:"ca"`
	KeyPath  string             `json:"key" yaml:"key"`
	CertPath string             `json:"cert" yaml:"cert"`
	Verify   bool               `json:"verify" yaml:"verify"`
	AuthType tls.ClientAuthType `json:"auth_type" yaml:"auth_type"`
}

func (s *SSL) Load() (*tls.Config, error) {
	var c tls.Config

	if len(s.CaPath) > 0 {
		caCert, err := ioutil.ReadFile(s.CaPath)
		if err != nil {
			return nil, fmt.Errorf("natsmq-ssql: Unable to load CA cert: %s", err)
		}
		caCertPool := x509.NewCertPool()

		caCertPool.AppendCertsFromPEM(caCert)

		c.ClientCAs = caCertPool
		c.InsecureSkipVerify = s.Verify
	}

	if len(s.CertPath) > 0 && len(s.KeyPath) > 0 {
		cert, err := tls.LoadX509KeyPair(s.CertPath, s.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("natsmq-ssl: Unable to load tls keypair: %s", err)
		}
		c.Certificates = append(c.Certificates, cert)
	}

	return &c, nil
}
