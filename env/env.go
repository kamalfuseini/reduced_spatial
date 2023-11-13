package env

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/kfuseini/reduced_spatial/config"
	"github.com/pkg/errors"
)

type Env struct {
	Config *config.Config
	Cass   *gocql.Session
}

func NewEnv(c *config.Config) (*Env, error) {
	e := &Env{
		Config: c,
	}

	cluster := gocql.NewCluster(c.Cassandra()...)
	cluster.Keyspace = "reduced_spatial"
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: c.CassUser,
		Password: c.CassPass,
	}
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{ Max: 30 * time.Second }
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.NumConns = 3
	cluster.Timeout = 10 * time.Second
	if cass, err := cluster.CreateSession(); err == nil {
		e.Cass = cass
	} else {
		return nil, errors.Wrap(err, "Failed to create Cassandra session")
	}

	return e, nil
}
