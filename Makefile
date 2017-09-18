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
	sudo docker build -t caops .


docker_run:
	sudo docker run --name caops01 -d -p 19042:9042 -p 18080:8080 -p 18787:8787 caops:latest
	export SEED_IP=`sudo docker inspect --format='{{ .NetworkSettings.IPAddress }}' caops01` ; \
	sudo docker run --name caops02 -d -e CASSANDRA_SEEDS="$$SEED_IP" -p 29042:9042 -p 28080:8080 -p 28787:8787 caops:latest ; \
	sudo docker run --name caops03 -d -e CASSANDRA_SEEDS="$$SEED_IP" -p 39042:9042 -p 38080:8080 -p 38787:8787 caops:latest


docker_stop:
	sudo docker ps | grep caops | awk '{print$$1}' | paste -s | xargs sudo docker stop

docker_prune:
	sudo docker rmi caops -f
	sudo docker container prune -f

docker_shell:
	sudo docker exec -it `sudo docker ps | grep caops01 | awk '{print$$1}'` /bin/bash
