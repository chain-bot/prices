#!/bin/bash

# Requires cc-test-reporter: https://docs.codeclimate.com/docs/configuring-test-coverage#section-locations-of-pre-built-binaries
source .env
cc-test-reporter before-build
cd app || exit
go test -coverprofile ../c.out $(go list ./... | grep -v /psql/generated)
cd ../
cc-test-reporter format-coverage -t gocov --prefix github.com/mochahub/coinprice-scraper c.out
cc-test-reporter after-build --prefix github.com/mochahub/coinprice-scraper -r=$CC_TEST_REPORTER_ID

mv c.out coverage/c.out