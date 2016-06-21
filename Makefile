SOURCEDIR := .
SOURCES := $(shell find $(SOURCEDIR) -name "*.go" -type f)

BINARYDIR := bin
BINARY := webmock-proxy

LDFLAGS = -ldflags="-w"

GLIDE := glide
GLIDE_VERSION := 0.10.2

.DEFAULT_GOAL := $(BINARYDIR)/$(BINARY)

$(BINARYDIR)/$(GLIDE):
	if [ ! -d $(BINARYDIR) ]; then mkdir $(BINARYDIR); fi
ifeq ($(shell uname),Darwin)
	curl -fL https://github.com/Masterminds/glide/releases/download/$(GLIDE_VERSION)/glide-$(GLIDE_VERSION)-darwin-amd64.zip -o glide.zip
	unzip glide.zip
	mv ./darwin-amd64/glide $(BINARYDIR)/$(GLIDE)
	rm -fr ./darwin-amd64
	rm ./glide.zip
else
	curl -fL https://github.com/Masterminds/glide/releases/download/$(GLIDE_VERSION)/glide-$(GLIDE_VERSION)-linux-amd64.zip -o glide.zip
	unzip glide.zip
	mv ./linux-amd64/glide $(BINARYDIR)/$(GLIDE)
	rm -fr ./linux-amd64
	rm ./glide.zip
endif

$(BINARYDIR)/$(BINARY): $(SOURCES)
	go build $(LDFLAGS) -o $(BINARYDIR)/$(BINARY)

deps: $(BINARYDIR)/$(GLIDE)
	$(BINARYDIR)/$(GLIDE) install
