export GO111MODULE=on

all: fp-multiuser

fp-multiuser:
	go build -o ./bin/fp-multiuser ./cmd/fp-multiuser
