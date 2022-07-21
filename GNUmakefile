TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=pagerduty

default: build

build: fmtcheck
	go install -tags timetzdata

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"


test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

local-install: build
	goreleaser build --single-target --snapshot --rm-dist
	set -e ;\
	VERSION=$$(jq -r '.version' dist/metadata.json) ;\
	OS=$$(jq -r '.runtime.goos' dist/metadata.json) ;\
	ARCH=$$(jq -r '.runtime.goarch' dist/metadata.json) ;\
	BIN_PATH=$$(jq -r '.[0].path' dist/artifacts.json) ;\
	BIN_ID=$$(jq -r '.[0].extra.ID' dist/artifacts.json) ;\
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/${PKG_NAME}/${PKG_NAME}/$${VERSION}/$${OS}_$${ARCH} ;\
	mv $${BIN_PATH} ~/.terraform.d/plugins/registry.terraform.io/${PKG_NAME}/${PKG_NAME}/$${VERSION}/$${OS}_$${ARCH}/$${BIN_ID} ;\
	echo "installed as version $${VERSION}" ;\

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test local-install

