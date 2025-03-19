package config

import (
	"flag"
	"log"
)

type Config struct {
	AdminKey             string
	PostresConnectionURL string
}

func MustLoad() Config {
	adminKeyString := flag.String(
		"admin-key",
		"",
		"pass-key to access admin features",
	)

	postgresConnectionString := flag.String(
		"postgres-url",
		"",
		"url string for postres",
	)

	flag.Parse()

	if *adminKeyString == "" {
		log.Fatal("admin key is not set")
	}
	if *postgresConnectionString == "" {
		log.Fatal("postres connection string is not specified")
	}

	return Config{
		AdminKey:             *adminKeyString,
		PostresConnectionURL: *postgresConnectionString,
	}
}
