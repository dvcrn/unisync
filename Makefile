apps/mackup:
	git submodule update --init --recursive
	go run cmd/importmackup/main.go

.PHONY: apps/mackup
