#!/bin/bash
echo $PWD
source ./.env
cd scraper || exit
godotenv -f ./.env gopherbadger
rm coverage.out
mv coverage_badge.png ../coverage_badge.png

#source ./.env
#godotenv -f ./.env go test $(go list ./... | grep -v /vendor/) -v -coverprofile cover.out . fmt
#go tool cover -html=cover.out -o cover.html
#gobadge -filename=cover.out
#if [[ "$OSTYPE" == "linux-gnu"* ]]; then
#  xdg-open cover.html
#elif [[ "$OSTYPE" == "darwin"* ]]; then
#  open cover.html
#fi


