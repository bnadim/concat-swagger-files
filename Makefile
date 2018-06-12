.PHONY: build

########################################################################################################################
#
# Parameters
#
########################################################################################################################
VERSION ?= 0.1

########################################################################################################################
#
# Do not edit
#
########################################################################################################################

########################################################################################################################
#
# HELP
#
########################################################################################################################

#COLORS
RED    := $(shell tput -Txterm setaf 1)
GREEN  := $(shell tput -Txterm setaf 2)
WHITE  := $(shell tput -Txterm setaf 7)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

# Add the following 'help' target to your Makefile
# And add help text after each target name starting with '\#\#'
# A category can be added with @category
HELP_HELPER = \
    %help; \
    while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^([a-zA-Z\-\%]+)\s*:.*\#\#(?:@([a-zA-Z\-\%]+))?\s(.*)$$/ }; \
    print "usage: make [target]\n\n"; \
    for (sort keys %help) { \
    print "${WHITE}$$_:${RESET}\n"; \
    for (@{$$help{$$_}}) { \
    $$sep = " " x (32 - length $$_->[0]); \
    print "  ${YELLOW}$$_->[0]${RESET}$$sep${GREEN}$$_->[1]${RESET}\n"; \
    }; \
    print "\n"; }

help: ##prints help
	@perl -e '$(HELP_HELPER)' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

build: ##@build Build for supported operating systems
	@rm -R build
	@env GOOS=darwin GOARCH=amd64 go build -o build/csf-${VERSION}-osx64/csf
	@env GOOS=linux GOARCH=amd64 go build -o build/csf-${VERSION}-linux64/csf
	@env GOOS=windows GOARCH=amd64 go build -o build/csf-${VERSION}-windows64/csf
	@env GOOS=windows GOARCH=386 go build -o build/csf-${VERSION}-windows32/csf

	@cd build/csf-${VERSION}-osx64 && zip -r ../csf-${VERSION}-osx64.zip csf
	@cd build/csf-${VERSION}-linux64 && zip -r ../csf-${VERSION}-linux64.zip csf
	@cd build/csf-${VERSION}-windows64 && zip -r ../csf-${VERSION}-windows64.zip csf
	@cd build/csf-${VERSION}-windows32 && zip -r ../csf-${VERSION}-windows32.zip csf