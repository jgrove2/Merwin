run: build
	@./bin/merwin_cli >> bin/logs/logfile

build:
	@go build -o bin/merwin_cli cmd/merwin/main.go

test:
	@./bin/merwin_cli