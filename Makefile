GO_COMPILER=go

ifeq ($(OS),Windows_NT)
	GO_COMPILER=go.exe
endif

all:
	@$(GO_COMPILER) build cmd/server/server.go && $(GO_COMPILER) build cmd/client/client.go
