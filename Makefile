BIN := ./bin/cutlink
MAIN_DB_FILE := database.db
SESSIONS_DB_FILE := sessions.db
MAIN_SRC := ./cmd


$(BIN): $(MAIN_DB_FILE) $(wildcard ./cmd/*.go) $(wildcard ./models/*.go)
	@mkdir -p ./bin
	CGO_ENABLED=1 go build -o $(BIN) $(MAIN_SRC)/...

run: $(MAIN_DB_FILE) $(wildcard ./cmd/*.go) $(wildcard ./db/*.go)
	CGO_ENABLED=1 go run $(MAIN_SRC)/...

$(MAIN_DB_FILE): sqlite.sql
	cat sqlite.sql | sqlite3 $(MAIN_DB_FILE)

clean:
	rm -rf ./bin $(SESSIONS_DB_FILE) $(MAIN_DB_FILE)

# Build docker image
docker: ./Dockerfile
	docker build -t cutlink .

# Build executable file in golang docker container
# Using bullseye version because of glibc backward compatibility
docker_exe: clean
	docker run --rm -v $(shell pwd):/app -w /app golang:1.21-bullseye make
