DB_URL ?= "./test.db"

folder ?= "internal"

all: install run

install:
	go install github.com/go-task/task/v3/cmd/task@latest && \
	go install golang.org/x/tools/gopls@latest && \
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
	go install github.com/marco-souza/hooker@latest && hooker init && \
	curl --proto '=https' --tlsv1.2 -LsSf https://github.com/frectonz/sql-studio/releases/download/0.1.16/sql-studio-installer.sh | sh

run: cmd/server/main.go
	task dev

deploy: ./fly.toml
	pkgx fly deploy --now -y

generate: sqlc.yml
	sqlc generate

studio:
	sql-studio sqlite ${DB_URL}

release: cmd/server/main.go generate
	CGO_CFLAGS="-D_LARGEFILE64_SOURCE" CGO_ENABLED=1 \
	go build -ldflags "-s -w" -o ./build/server ./cmd/server/main.go

fmt:
	go fmt ./... && npx prettier -w views ./README.md ./docker-compose.yml

t: test
	task tests

test: ./tests/
	go test -v ./...

encrypt: .env
	gpg -c .env

decrypt: .env.gpg
	gpg -d .env.gpg > .env

gen:
	@go run ./cmd/cli/cli.go ${folder} ${name}
