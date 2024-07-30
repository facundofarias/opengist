.PHONY: go_mod

go_mod:
	@if [ ! -f go.mod ]; then \
		go mod init <module-name>; \
	fi
	go mod tidy
	go mod download
	go mod verify
