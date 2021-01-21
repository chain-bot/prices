#!/bin/bash
echo $PWD

source ./.env
echo $psql_user
sqlboiler psql
