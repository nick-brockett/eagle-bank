package postgres

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	Username     string `env:"POSTGRES_USER, default=postgres"`
	Password     string `env:"POSTGRES_PASSWORD, default=password123"`
	Host         string `env:"POSTGRES_HOST, default=localhost"` //default=postgres-db"`
	Port         string `env:"POSTGRES_PORT, default=5432"`
	DatabaseName string `env:"POSTGRES_DB, default=postgres"`
	SSLMode      string `env:"POSTGRES_SSL_MODE, default=disable"`
}

func (c Config) GetUsername() string {
	return c.Username
}

func (c Config) GetPassword() string {
	return c.Password
}

func (c Config) GetHost() string {
	return c.Host
}

func (c Config) GetPort() int {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return 0
	}
	return port
}

func (c Config) GetDatabaseName() string {
	return c.DatabaseName
}

func (c Config) PostgresConnString() (string, error) {
	if isStrEmpty(c.Host) {
		return "", errors.New("host is required")
	}
	if isStrEmpty(c.Port) {
		return "", errors.New("port is required")
	}
	if isStrEmpty(c.Username) {
		return "", errors.New("username is required")
	}
	if isStrEmpty(c.Password) {
		return "", errors.New("password is required")
	}
	if isStrEmpty(c.DatabaseName) {
		return "", errors.New("database name is required")
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.GetHost(), c.GetPort(), c.GetUsername(), c.GetPassword(), c.GetDatabaseName()), nil
}
func isStrEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
