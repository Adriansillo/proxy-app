package utils

import (
	env "github.com/joho/godotenv"
)

func LoadEnv() {
	env.Load()
}
