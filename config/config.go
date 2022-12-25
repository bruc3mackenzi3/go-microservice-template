package config

import "os"

const Port = 8080
const Test = 1234

var Env Environment

type Environment struct {
	PostgresHost     string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
}

func init() {
	Env = Environment{
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
	}
}
