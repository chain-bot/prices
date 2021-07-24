package configs

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// LoadEnv loads env vars from .env at root of repo
func GetProjectRoot() string {
	re := regexp.MustCompile(`^(.*` + PROJECT_DIR + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	return string(rootPath)
}

func GetMigrationDir() string {
	return fmt.Sprintf("file://%s%s", GetProjectRoot(), MIGRATION_DIR)
}
func LoadEnv() {
	env := Environment(os.Getenv("CHAINBOT_ENV"))
	if env != LocalEnv && env != NilEnv {
		// Env variables already set
		return
	}
	re := regexp.MustCompile(`^(.*` + PROJECT_DIR + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	filePath := fmt.Sprintf("%s/.env", string(rootPath))
	err := godotenv.Load(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"cause":    err,
			"cwd":      cwd,
			"filePath": filePath,
		}).Fatal("problem loading .env file")
		os.Exit(-1)
	}
}
