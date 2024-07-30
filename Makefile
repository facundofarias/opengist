.PHONY: go_mod check_changes

go_mod:
	@if [ ! -f go.mod ]; then \
		go mod init github.com/thomiceli/opengist; \
	fi
	go mod tidy

check_changes: go_mod
	go mod download
	go mod verify
