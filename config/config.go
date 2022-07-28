package config

import "os"

const Port = 80

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
