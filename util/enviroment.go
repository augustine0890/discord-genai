package util

import "github.com/joho/godotenv"

func LoadEnv(env string) error {
	switch env {
	case "prod":
		err := godotenv.Load("prod.env")
		return err
	case "dev":
		err := godotenv.Load("dev.env")
		return err
	}
	return nil
}
