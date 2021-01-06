source ./.env
mkdir -p models
xo pgsql://$POSTGRES_USERNAME:$POSTGRES_PASSWORD@$POSTGRESQL_HOST:5432/$POSTGRES_DATABASE?sslmode=disable/ -v -s last_sync -o models
go build ./models/
go install ./models/