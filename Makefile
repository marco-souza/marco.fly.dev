DB_URL ?= "./test.db"

folder ?= "internal"
count ?= 1
time ?= "1s"
test ?= "."

all: install run

install:
	go install github.com/go-task/task/v3/cmd/task@latest && \
	go install golang.org/x/tools/gopls@latest && \
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
	go install github.com/marco-souza/hooker@latest && hooker init && \
	curl --proto '=https' --tlsv1.2 -LsSf https://github.com/frectonz/sql-studio/releases/download/0.1.16/sql-studio-installer.sh | sh

run: cmd/server/main.go
	task dev

deploy: secrets ./fly.toml
	pkgx fly deploy --now -y

secrets: .env
	for line in $$(cat .env | grep -v '^#'); do \
		echo $$line | sed -e 's#=.*##g'; \
		echo $$line | pkgx fly secrets set $$line --stage ; \
	done

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

b: bench

bench: ./tests/
	go test -bench=${test} ./tests/bench/... -count=${count} -benchmem -benchtime=${time}

encrypt: .env
	gpg -c .env
	sed -e 's/=.*/=""/g' .env > .env.example

decrypt: .env.gpg
	gpg -d .env.gpg > .env

gen:
	@go run ./cmd/cli/cli.go ${folder} ${name}
