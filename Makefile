APPNAME     := athena
VERSION     := $(shell git describe --all --always --dirty --long)
LDFLAGS     := "-X bitbucket.org/crossengage/athena/cmd/athena.version=$(VERSION) -X bitbucket.org/crossengage/athena/cmd/athena.appName=$(APPNAME)"
PLATFORMS   := darwin-386 darwin-amd64 linux-386 linux-amd64 linux-arm windows-386 windows-amd64
INSTALL_PKG := ./cmd/athena
BIN_DIR     := ./bin
DIST_DIR    := ./dist
APP_BIN     := $(BIN_DIR)/$(APPNAME)
GO_FILES     = $(shell find ./ -type f -name '*.go')


default: build

help:
	@echo "make [target]"
	@echo " Targets: "
	@echo "   build     Build for current OS+Arch into $(BIN_DIR) "
	@echo "   dist      Build for many platforms into $(DIST_DIR) "
	@echo "   clean     Cleanup binaries                          "
	@echo "   deps      Install deps for the project              "


.PHONY: clean
clean:
	go clean
	rm -fv $(APP_BIN)
	rm -fv $(DIST_DIR)/*


.PHONY: deps
deps:
	go get -v .


build: $(APP_BIN)
$(APP_BIN): $(GO_FILES)
	mkdir -p $(BIN_DIR)
	go build -ldflags=$(LDFLAGS) -o $@ $(INSTALL_PKG)


dist: $(PLATFORMS)
$(PLATFORMS):
	$(eval GOOS := $(firstword $(subst -, ,$@)))
	$(eval GOARCH := $(lastword $(subst -, ,$@)))
	mkdir -p $(DIST_DIR)
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(LDFLAGS) -o $(DIST_DIR)/$(APPNAME).$@ $(INSTALL_PKG)