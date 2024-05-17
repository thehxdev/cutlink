# Using bullseye version because of glibc backward compatibility
FROM golang:1.22-bullseye

WORKDIR /app
COPY . .
EXPOSE 5000

RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install sqlite3 -y && apt-get install build-essential -y
RUN mkdir ./bin
# RUN cat sqlite.sql | sqlite3 database.db
# ENV GOPROXY=direct
ENV CGO_ENABLED=1
RUN go build -o ./bin/cutlink ./cmd/...

CMD ["./bin/cutlink"]
