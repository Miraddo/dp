SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

run:
	go run .

test:
	go test .
	
tidy:
	go mod tidy
	go mod vendor