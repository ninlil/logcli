build:
	@mkdir -p ./bin
	go build -o bin/logcli .

install: build
	@mkdir -p ~/bin
	install bin/logcli ~/bin/logcli