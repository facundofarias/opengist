.PHONY: go_mod

go_mod:
	@if [ ! -f go.mod ]; then \
		go mod init github.com/thomiceli/opengist; \
	fi
	go mod tidy
	go mod download
	go mod verify
