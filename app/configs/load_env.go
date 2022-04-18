package configs

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// GetProjectRoot loads env vars from .env at root of repo
func GetProjectRoot() string {
	re := regexp.MustCompile(`^(.*` + projectDir + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	return string(rootPath)
}

func GetMigrationDir() string {
	return fmt.Sprintf("file://%s%s", GetProjectRoot(), migrationDir)
}
func LoadEnv() {
	env := Environment(os.Getenv("CHAINBOT_ENV"))
	if env != LocalEnv && env != NilEnv {
		log.WithField("env", env).Debug("skipping .env file")
		return
	}
	rootPath := GetProjectRoot()
	filePath := fmt.Sprintf("%s/.env", rootPath)
	err := godotenv.Load(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"cause":    err,
			"filePath": filePath,
		}).Fatal("problem loading .env file")
		os.Exit(-1)
	}
}
