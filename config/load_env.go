package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
)

// LoadEnv loads env vars from .env at root of repo
func GetProjectRoot() string {
	re := regexp.MustCompile(`^(.*` + PROJECT_DIR + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	return string(rootPath)
}
func LoadEnv() {
	if os.Getenv("COINPRICE_ENV") != "" {
		// Env variables already set
		return
	}
	re := regexp.MustCompile(`^(.*` + PROJECT_DIR + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.WithFields(log.Fields{
			"cause": err,
			"cwd":   cwd,
		}).Fatal("Problem loading .env file")

		os.Exit(-1)
	}
}
