#!/bin/bash

# Requires cc-test-reporter: https://docs.codeclimate.com/docs/configuring-test-coverage#section-locations-of-pre-built-binaries
# source .env
./cc-test-reporter before-build
cd app || exit
go test -coverprofile ../c.out $(go list ./... | grep -v /psql/generated)
cd ../
./cc-test-reporter format-coverage -t gocov --prefix github.com/chain-bot/scraper c.out
./cc-test-reporter after-build --prefix github.com/chain-bot/scraper -r=$CC_TEST_REPORTER_ID
go tool cover -html=c.out -o=coverage/c.html
mv c.out coverage/c.out