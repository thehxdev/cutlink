BUILD_DIR := ./build
BIN := $(BUILD_DIR)/cutlink
MAIN_DB_FILE := database.db
SESSIONS_DB_FILE := sessions.db
MAIN_SRC := ./cmd

all: $(BIN)

$(BIN): $(wildcard ./cmd/*.go) $(wildcard ./models/*.go) $(wildcard ./rand/*.go)
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -o $(BIN) $(MAIN_SRC)/...

run: $(wildcard ./cmd/*.go) $(wildcard ./db/*.go)
	CGO_ENABLED=1 go run $(MAIN_SRC)/...

# $(MAIN_DB_FILE): sqlite.sql
# 	cat sqlite.sql | sqlite3 $(MAIN_DB_FILE)

# make database.db file
# db: $(MAIN_DB_FILE)

clean:
	rm -rf $(BUILD_DIR) $(SESSIONS_DB_FILE) $(MAIN_DB_FILE)
	# go clean

# Build docker image
docker: clean ./Dockerfile
	docker build -t cutlink .

# Build executable file in golang docker container
# Using bullseye version because of glibc backward compatibility
docker_exe:
	docker run --rm -v $(shell pwd):/app -w /app golang:1.22-bullseye make
