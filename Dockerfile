# base stage
FROM golang:1.22-alpine as base
WORKDIR /app
COPY ./views/ ./views/
COPY ./static/ ./static/
COPY ./scripts/ ./scripts/

# pre-build stage
FROM base as pre-build

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download
RUN apk add --no-cache make build-base

COPY . .

# dev stage
FROM pre-build as dev
RUN make install
CMD ["make", "release"]

# build stage
FROM pre-build as build
RUN make release

# prod stage
FROM base as prod
COPY --from=build /app/build/server ./
CMD ["/app/server"]
