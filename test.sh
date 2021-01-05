source ./.env
#go test $(go list ./... | grep -v /vendor/) -v -coverprofile cover.out . fmt
#go tool cover -html=cover.out -o cover.html
#
#if [[ "$OSTYPE" == "linux-gnu"* ]]; then
#  xdg-open cover.html
#elif [[ "$OSTYPE" == "darwin"* ]]; then
#  open cover.html
#fi
gopherbadger -md="README.md"

