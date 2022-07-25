package config

import "os"

var Env Environment

type Environment struct {
	PostgresHost     string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
}

func init() {
	Env = Environment{
		PostgresHost:     os.Getenv("PGHOST"),
		PostgresDB:       os.Getenv("PGDB"),
		PostgresUser:     os.Getenv("PGUSER"),
		PostgresPassword: os.Getenv("PGPASSWORD"),
	}
}
