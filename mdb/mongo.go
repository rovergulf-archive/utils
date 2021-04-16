package mdb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

const (
	srvUriPrefix     = "mongodb+srv"
	clusterUriPrefix = "mongodb"
)

type MongoDB struct {
	Client *mongo.Client
	lg     *zap.SugaredLogger
	cfg    *Config
}

type Config struct {
	Addr     string             `json:"addr" yaml:"addr"`
	User     string             `json:"user" yaml:"user"`
	Password string             `json:"password" yaml:"password"`
	Database string             `json:"database" yaml:"database"`
	Provider string             `json:"provider" yaml:"provider"`
	Logger   *zap.SugaredLogger `json:"-" yaml:"-"`
}

func uriPrefix(pv string) string {
	switch pv {
	case "cluster":
		return clusterUriPrefix
	default:
		return srvUriPrefix
	}
}

func (c *Config) Uri() string {
	if len(c.User) > 0 && len(c.Password) > 0 {
		return fmt.Sprintf("%s://%s:%s@%s", uriPrefix(c.Provider), c.User, c.Password, c.Addr)
	}
	return fmt.Sprintf("%s://%s", uriPrefix(c.Provider), c.Addr)
}

func NewClient(ctx context.Context, lg *zap.SugaredLogger, c *Config) (*MongoDB, error) {
	if c.Logger == nil {
		c.Logger = lg.Named("mdb")
	}

	if c.Provider == "" {
		c.Provider = srvUriPrefix
	}

	mdb := &MongoDB{
		lg:  c.Logger.Named("mongodb"),
		cfg: c,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.Uri()))
	if err != nil {
		mdb.lg.Errorf("Unable to connect mongo db instance: %s", err)
		return nil, err
	}

	mdb.Client = client
	if err := mdb.Client.Ping(ctx, readpref.Primary()); err != nil {
		mdb.lg.Errorf("Unable to ping mongo connection: %s", err)
		return nil, err
	}

	mdb.lg.Debugw("Connected to MongoDB instance", "addr", c.Addr, "user", c.User)

	return mdb, nil
}

func (m *MongoDB) Shutdown(ctx context.Context) {
	if err := m.Client.Disconnect(ctx); err != nil {
		m.lg.Errorf("Unable to close MongoDB connection: %s", err)
	}
}

func (m *MongoDB) GetCollection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return m.Client.Database(m.cfg.Database).Collection(name, opts...)
}
