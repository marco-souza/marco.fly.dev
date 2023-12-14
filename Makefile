deps:
	go install github.com/cosmtrek/air@latest && \
	go install golang.org/x/tools/gopls@latest

run: deps cmd/main.go
	air

deploy: ./fly.toml
	pkgx fly deploy --now -y

build: cmd/main.go
	CGO_ENABLED=1 go build -o ./build/server ./cmd/main.go

fmt:
	go fmt ./... && bunx prettier -w views ./README.md ./docker-compose.yml

t: test
test: ./tests/
	go test -v ./tests/

encrypt: .env
	gpg -c .env

decrypt: .env.gpg
	gpg -d .env.gpg > .env
