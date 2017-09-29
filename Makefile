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
	for i in `seq 1 3`; do vagrant up $(DEV_CASS_VER)x0$$i; done


provision: build
	for i in `seq 1 3`; do vagrant up $(DEV_CASS_VER)x0$$i --provision; done

	
docker:
	docker build -t caops .


docker_run:
	docker run --name caops01 -d -e CASSANDRA_BROADCAST_ADDRESS=172.17.255.255 -p 19042:9042 -p 18080:8080 -p 18787:8787 caops:latest
	export SEED_IP=`docker inspect --format='{{ .NetworkSettings.IPAddress }}' caops01` ; \
	docker run --name caops02 -d -e CASSANDRA_SEEDS=$$SEED_IP -p 29042:9042 -p 28080:8080 -p 28787:8787 caops:latest


docker_stop:
	docker ps | grep caops | awk '{print$$1}' | paste -s | xargs docker stop
	docker rmi caops -f
	docker container prune -f


docker_shell_01:
	docker exec -it `docker ps | grep caops01 | awk '{print$$1}'` /bin/bash

docker_shell_02:
	docker exec -it `docker ps | grep caops02 | awk '{print$$1}'` /bin/bash
