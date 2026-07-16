package utils

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var DB_NAME string
var DB_PORT string
var DB_HOST string
var DB_USER string
var DB_PASS string

var SERVER_PORT string
var SERVER_HOST string

var ADMIN_NAME string
var ADMIN_EMAIL string
var ADMIN_PASS string

var JWT_SECRET string
var JWT_EXPIRY string

var SECURE_TOKEN_LENGTH int
var SECURE_TOKEN_EXPIRY string

func InitEnv() error {
	var err = godotenv.Load(".env")
	if err != nil {
		return err
	}

	DB_NAME = os.Getenv("DB_NAME")
	DB_PORT = os.Getenv("DB_PORT")
	DB_HOST = os.Getenv("DB_HOST")
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")

	SERVER_PORT = os.Getenv("SERVER_PORT")
	SERVER_HOST = os.Getenv("SERVER_HOST")

	ADMIN_NAME = os.Getenv("ADMIN_NAME")
	ADMIN_EMAIL = os.Getenv("ADMIN_EMAIL")
	ADMIN_PASS = os.Getenv("ADMIN_PASS")

	JWT_SECRET = os.Getenv("JWT_SECRET")
	JWT_EXPIRY = os.Getenv("JWT_EXPIRY")

	SECURE_TOKEN_LENGTH, err = strconv.Atoi(os.Getenv("SECURE_TOKEN_LENGTH"))
	if err != nil {
		return err
	}
	SECURE_TOKEN_EXPIRY = os.Getenv("SECURE_TOKEN_EXPIRY")

	return nil
}
