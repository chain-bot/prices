package configs

const (
	PROJECT_DIR   = "prices"
	MIGRATION_DIR = "/app/internal/data/psql/migrations"
)

const (
	NilEnv   = Environment("")
	LocalEnv = Environment("LOCAL")
	DevEnv   = Environment("DEV")
	ProdEnv  = Environment("ENV")
)
