SHELL=bash
OK_MSG = \x1b[32m ✔\x1b[0m
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOLIST?=$$(go list ./... | grep -v vendor)

default: test

integration:
	go test -v ./tests/integration

# tools fetches necessary dev requirements
tools:
	go get -u github.com/robertkrimen/godocdown/godocdown
	go get -u github.com/kardianos/govendor
	go get -u honnef.co/go/tools/cmd/gosimple
	go get -u honnef.co/go/tools/cmd/unused
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u github.com/client9/misspell/cmd/misspell
	go get -u github.com/golang/lint/golint

vendor-status:
	@govendor status

coverprofile:
	@go test ./pagerduty/... -coverprofile coverage.out && go tool cover -html=coverage.out

lint:
	@echo -n "==> Checking that code complies with golint requirements..."
	@ret=0 && for pkg in $(GOLIST); do \
		test -z "$$(golint $$pkg | tee /dev/stderr)" || ret=1; \
		done ; exit $$ret
	@echo -e "$(OK_MSG)"

# check combines all checks into a single command
check: fmtcheck vet misspell staticcheck simple unused lint vendor-status

# fmt formats Go code.
fmt:
	gofmt -w $(GOFMT_FILES)

test: check
	@echo "==> Checking that code complies with unit tests..."
	@go test $(GOLIST) -cover

webdoc:
	@echo "==> Starting webserver at http://localhost:6060"
	@sleep 1 && open http://localhost:6060 &
	@godoc -http=:6060

unused:
	@echo -n "==> Checking that code complies with unused requirements..."
	@unused $(GOLIST)
	@echo -e "$(OK_MSG)"

fmtcheck:
	@echo -n "==> Checking that code complies with gofmt requirements..."
	@gofmt_files=$$(gofmt -l $(GOFMT_FILES)) ; if [[ -n "$$gofmt_files" ]]; then \
		echo 'gofmt needs running on the following files:'; \
		echo "$$gofmt_files"; \
		echo "You can use the command: \`make fmt\` to reformat code."; \
		exit 1; \
	fi
	@echo -e "$(OK_MSG)"

misspell:
	@echo -n "==> Checking for misspelling errors..."
	@misspell --error $(GOFMT_FILES)
	@echo -e "$(OK_MSG)"

simple:
	@echo -n "==> Checking that code complies with gosimple requirements..."
	@gosimple $(GOLIST)
	@echo -e "$(OK_MSG)"

staticcheck:
	@echo -n "==> Checking that code complies with staticcheck requirements..."
	@staticcheck $(GOLIST)
	@echo -e "$(OK_MSG)"

vet:
	@echo -n "==> Checking that code complies with go vet requirements..."
	@go vet $(GOLIST) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
	@echo -e "$(OK_MSG)"
