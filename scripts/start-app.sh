#!/bin/bash

# Run from Root of repo
go run app/cmd/server/main.go &
go run app/cmd/scraper/main.go

