package natsmq

import (
	"crypto/tls"
	"github.com/nats-io/nats.go"
)

// Config represent stan.Conn and nats.Conn connection parameters
type Config struct {
	ClusterId string            `json:"clusterId" yaml:"cluster_id"`
	BrokerId  string            `json:"brokerId" yaml:"broker_id"`
	Broker    string            `json:"brokers" yaml:"broker"`
	User      string            `json:"user" yaml:"user"`
	Password  string            `json:"password" yaml:"password"`
	Token     string            `json:"token" yaml:"token"`
	Channels  map[string]string `json:"channels" yaml:"channels"`
	TLS       SSL               `json:"tls" yaml:"tls"`
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

func (c *Config) GetNatsUserInfo() nats.Option {
	return nats.UserInfo(c.User, c.Password)
}
