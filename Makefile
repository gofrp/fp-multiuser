export GO111MODULE=on

all: frps-multiuser

frps-multiuser:
	go build -o ./bin/frps-multiuser ./cmd/fp-multiuser
