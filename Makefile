export GO111MODULE=on

copy:
	mkdir ./bin
	cp ./config/frps-multiuser.ini ./bin/frps-multiuser.ini
	cp -r ./assets/ ./bin/assets/

frps-multiuser: copy
	go build -o ./bin/frps-multiuser ./cmd/fp-multiuser
