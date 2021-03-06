APPNAME       := CaOps
VERSION       := $(shell git describe --all --always --dirty --long)
LDFLAGS       := "-X main.appName=$(APPNAME) -X main.version=$(VERSION)"
PLATFORMS     := darwin-amd64 linux-amd64 linux-arm windows-amd64
INSTALL_PKG   := ./cmd/CaOps
BIN_DIR       := ./bin
DIST_DIR      := ./dist
APP_BIN       := $(BIN_DIR)/$(APPNAME)
GO_FILES       = $(shell find ./ -type f -name '*.go')
DEV_CASS_VER  := c22


.PHONY: default
default: build


.PHONY: ignored
ignored:
	git ls-files --others -i --exclude-standard


.PHONY: help
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
	go build -race -ldflags=$(LDFLAGS) -o $@ $(INSTALL_PKG)


dist: $(PLATFORMS)
$(PLATFORMS):
	$(eval GOOS := $(firstword $(subst -, ,$@)))
	$(eval GOARCH := $(lastword $(subst -, ,$@)))
	mkdir -p $(DIST_DIR)
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(LDFLAGS) -o $(DIST_DIR)/$(APPNAME).$@ $(INSTALL_PKG)


dev_vms: build
	for i in `seq 1 2`; do vagrant up $(DEV_CASS_VER)x0$$i; done


provision: build
	for i in `seq 1 2`; do vagrant up --provision $(DEV_CASS_VER)x0$$i; done
