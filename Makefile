export GO111MODULE=on

build: frps-multiuser
	cp ./config/frps-multiuser.ini ./bin/frps-multiuser.ini
	cp -r ./assets/ ./bin/assets/

frps-multiuser:
	go build -o ./bin/frps-multiuser ./cmd/fp-multiuser
