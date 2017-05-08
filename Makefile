BIN := athena
LDFLAGS := "-w"

help:
	@echo "make [target]"
	@echo " Targets: "
	@echo "   linux       Build for Linux amd64 into dist/   "
	@echo "   darwin      Build for Darwin amd64 into dist/  "
	@echo "   windows     Build for Windows amd64 into dist/ "
	@echo "   dist        Build for all above into dist/     "


clean:
	go clean


deps:
	go get -v .


VERSION := $(shell git describe --all --always --dirty --long)
.PHONY: version
version:
	@printf "package cmd\nconst version=\`$(VERSION)\`\n" | \
		gofmt | tee cmd/gen_version.go


build: version
	go generate
	go build -v


linux: version
	$(info Building for Linux 64bits...)
	go generate && \
	env GOOS=linux GOHOSTARCH=amd64 go build -ldflags="$(LDFLAGS)" -v -o dist/$(VERSION)/linux-amd64/$(BIN)


darwin: version
	$(info Building for Darwin 64bits...)
	go generate && \
	env GOOS=darwin GOHOSTARCH=amd64 go build -ldflags="$(LDFLAGS)" -v -o dist/$(VERSION)/darwin-amd64/$(BIN)


windows: version
	$(info Building for Windows 64bits...)
	go generate && \
	env GOOS=windows GOHOSTARCH=amd64 go build -ldflags="$(LDFLAGS)" -v -o dist/$(VERSION)/windows-amd64/$(BIN).exe


dist: linux darwin windows
	$(info Done.)