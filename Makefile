deps:
	go install github.com/cosmtrek/air@latest && \
	go install github.com/a-h/templ/cmd/templ@latest

run: gen cmd/main.go
	air

deploy: ./fly.toml
	pkgx fly deploy --now -y

build: cmd/main.go
	go build -o ./build/server ./cmd/main.go

gen: deps
	templ generate

fmt:
	go fmt ./... && \
  bunx prettier -w views ./README.md ./docker-compose.yml

t: ./tests/
	go test -v ./tests/

encrypt: .env
	gpg -c .env

decrypt: .env.gpg
	gpg -d .env.gpg > .env
