BIN := ./bin/cutlink
DB_FILE := database.db
MAIN_SRC := ./cmd

$(BIN): $(DB_FILE) $(wildcard ./cmd/*.go) $(wildcard ./models/*.go)
	@mkdir -p ./bin
	go build -o $(BIN) $(MAIN_SRC)/...

run: $(DB_FILE) $(wildcard ./cmd/*.go) $(wildcard ./db/*.go)
	go run $(MAIN_SRC)/...

$(DB_FILE): sqlite.sql
	cat sqlite.sql | sqlite3 $(DB_FILE)


clean:
	rm -rf ./bin $(DB_FILE)
