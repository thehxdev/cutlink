FROM golang:1.21.1-bookworm

WORKDIR /app
COPY . .
EXPOSE 5000

RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install sqlite3 -y && apt-get install build-essential -y
RUN mkdir ./bin
RUN cat sqlite.sql | sqlite3 database.db
RUN CGO_ENABLED=1 go build -o ./bin/cutlink ./cmd/...

ENTRYPOINT ["./bin/cutlink"]
