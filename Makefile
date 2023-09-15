cutlink: $(wildcard ./cmd/*.go) $(wildcard ./db/*.go)
	@mkdir -p ./bin
	go build -o ./bin/cutlink ./cmd/...

run: $(wildcard ./cmd/*.go) $(wildcard ./db/*.go)
	go run ./cmd/...

db: sqlite.sql
	cat sqlite.sql | sqlite3 database.db


clean:
	rm -rf ./bin ./database.db
