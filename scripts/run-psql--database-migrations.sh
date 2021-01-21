#!/bin/bash
echo $PWD
source ./.env

go run data/psql/migrations/main/run_migrations.go

