#!/bin/bash

cd app || exit
go test -short -coverprofile ../c.out $(go list ./... | grep -v /psql/generated) || exit $?
echo $?
cd ../

#source ./.env
#cd app || exit
##go test -v $(go list ./... | grep -v /psql/generated)
#go test -v $(go list ./... | grep -v /psql/generated) -v -coverprofile ../cover.out . fmt
#cd ../
#go tool cover -html=cover.out -o=cover.html
#go tool cover -func=cover.out -o=cover.out
#gobadge -filename=cover.out
#rm -rf cover.out
#mv cover.html docs/cover.html
#if [[ "$OSTYPE" == "linux-gnu"* ]]; then
#  xdg-open cover.html
#elif [[ "$OSTYPE" == "darwin"* ]]; then
#  open cover.html
#fi


