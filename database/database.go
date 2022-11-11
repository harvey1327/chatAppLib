package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	Close() error
	getDatabase() *mongo.Database
}

type mongoDBImpl struct {
	database *mongo.Database
	client   *mongo.Client
}

type dBConfig struct {
	host     string
	port     int
	username string
	password string
}

const USER = "user"

func DBConfig(host string, port int, username string, password string) dBConfig {
	return dBConfig{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (dbc dBConfig) validate() error {
	if dbc.host == "" {
		return fmt.Errorf("database host is invalid: '%s'", dbc.host)
	} else if dbc.port <= 0 {
		return fmt.Errorf("database port is invalid: '%d'", dbc.port)
	} else if dbc.username == "" {
		return fmt.Errorf("database username is invalid: '%s'", dbc.username)
	} else if dbc.password == "" {
		return fmt.Errorf("database password is invalid: '%s'", dbc.password)
	} else {
		return nil
	}
}

func NewDB(database string, config dBConfig) DB {
	err := config.validate()
	if err != nil {
		log.Fatal(err)
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", config.host, config.port)).SetAuth(options.Credential{Username: config.username, Password: config.password}))
	if err != nil {
		log.Fatal(err)
	}
	return &mongoDBImpl{
		database: client.Database(database),
		client:   client,
	}
}

func (m *mongoDBImpl) Close() error {
	return m.client.Disconnect(context.TODO())
}

func (m *mongoDBImpl) getDatabase() *mongo.Database {
	return m.database
}
