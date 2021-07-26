FROM golang:1.16-alpine
RUN apk add --no-cache git
RUN apk add --no-cache bash
# Set the Current Working Directory inside the container
#Note: Need to have "prices" because of how we do pathing in some parts of the code
WORKDIR /prices
# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY app ./app
# Build the Go app
RUN go build -o ./out-server ./app/cmd/server/main.go

ENV CHAINBOT_ENV=TEST
# Run the binary program produced by `go install`
CMD ["./out-server"]
