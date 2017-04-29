BIN := athena
LDFLAGS := "-w"

help:
	@echo "make [target]"
	@echo " Targets: "
	@echo "   dev-deps    Install development dependencies   "
	@echo "   vm          Initialize Virtualbox VM           "
	@echo "   migrate     Rollout migrations                 "
	@echo "   reset       Rollback and rollout migrations    "
	@echo "   down        Rollback migrations                "
	@echo "   linux       Build for Linux amd64 into dist/   "
	@echo "   darwin      Build for Darwin amd64 into dist/  "
	@echo "   windows     Build for Windows amd64 into dist/ "
	@echo "   dist        Build for all above into dist/     "


clean:
	go clean


deps:
	go get -v .


VERSION := $(shell git describe --all --always --dirty --long)
.PHONY: gen_version.go
gen_version.go:
	@printf "package main\nconst version=\`$(VERSION)\`\n" | \
		gofmt | tee $@


dev-deps:
	go get github.com/tools/godep
	go get github.com/mattes/migrate


vm:	dev-deps
	vagrant box update
	vagrant up


migrate: dev-deps
	migrate -url cassandra://$(VMIP)/cassandra_dump -path ./migrations up


reset: dev-deps
	migrate -url cassandra://$(VMIP)/cassandra_dump -path ./migrations reset


down: dev-deps
	migrate -url cassandra://$(VMIP)/cassandra_dump -path ./migrations down


build:
	go generate
	go build -v


linux:
	$(info Building for Linux 64bits...)
	go generate && \
	env GOOS=linux GOHOSTARCH=amd64 go build -ldflags="$(LDFLAGS)" -v -o dist/$(VERSION)/linux-amd64/$(BIN)


darwin:
	$(info Building for Darwin 64bits...)
	go generate && \
	env GOOS=darwin GOHOSTARCH=amd64 go build -ldflags="$(LDFLAGS)" -v -o dist/$(VERSION)/darwin-amd64/$(BIN)


windows:
	$(info Building for Windows 64bits...)
	go generate && \
	env GOOS=windows GOHOSTARCH=amd64 go build -ldflags="$(LDFLAGS)" -v -o dist/$(VERSION)/windows-amd64/$(BIN).exe


dist: linux darwin windows
	$(info Done.)