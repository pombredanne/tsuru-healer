# Copyright 2012 tsuru-healer authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

define HG_ERROR

FATAL: you need mercurial (hg) to download tsuru-healer dependencies.
       Check INSTALL.md for details
endef

define GIT_ERROR

FATAL: you need git to download tsuru-healer dependencies.
       Check INSTALL.md for details
endef

define BZR_ERROR

FATAL: you need bazaar (bzr) to download tsuru-healer dependencies.
       Check INSTALL.md for details
endef

all: check-path get test

# It does not support GOPATH with multiple paths.
check-path:
ifndef GOPATH
	@echo "FATAL: you must declare GOPATH environment variable, for more"
	@echo "       details, please check INSTALL.md file and/or"
	@echo "       http://golang.org/cmd/go/#GOPATH_environment_variable"
	@exit 1
endif
ifneq ($(subst ~,$(HOME),$(GOPATH))/src/github.com/globocom/tsuru-healer, $(PWD))
	@echo "FATAL: you must clone tsuru-healer inside your GOPATH To do so,"
	@echo "       you can run go get github.com/globocom/tsuru-healer/..."
	@echo "       or clone it manually to the dir $(GOPATH)/src/github.com/globocom/tsuru-healer"
	@exit 1
endif

get: hg git bzr get-test get-prod

hg:
	$(if $(shell hg), , $(error $(HG_ERROR)))

git:
	$(if $(shell git), , $(error $(GIT_ERROR)))

bzr:
	$(if $(shell bzr), , $(error $(BZR_ERROR)))

get-test:
	@/bin/echo -n "Installing test dependencies... "
	@go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | tr ' ' '\n' |\
		grep -v 'github.com/globocom/tsuru-healer' |\
		sort | uniq | xargs go get -u
	@/bin/echo "ok"

get-prod:
	@/bin/echo -n "Installing production dependencies... "
	@go list -f '{{range .Imports}}{{.}} {{end}}' ./... | tr ' ' '\n' |\
		grep -v 'github.com/globocom/tsuru-healer' |\
		sort | uniq | xargs go get -u
	@/bin/echo "ok"

test:
	@go test -i ./...
	@for pkg in `go list ./...`; do go test $$pkg; done
