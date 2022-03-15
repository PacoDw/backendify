GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TEST=./...
GOLANGCI_VERSION=latest

# test all the existing test files or just one
test: fmtcheck
	@if [[ "${pkg}" = "" && "${name}" = "" ]] ; then \
		go clean -testcache && go test $(TEST) -v -timeout=30s -parallel=4 ; \
	else \
		cd ./$(pkg) ; \
	    go clean -testcache && go test $(TEST) -v -run "$(name)" -timeout=30s -parallel=4 ; \
	fi

clean-cache:
	@go clean -cache -modcache -i -r

# check what files need to format
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

# check if the tools are installed
checktools:
	@sh -c "'$(CURDIR)/scripts/checktools.sh'"

# format all the files with gofmt
fmt:
	@if [ "${pkg}" = "" ] ; then \
		echo "==> Fixing source code with gofmt..." ; \
		gofmt -s -w ./ ; \
	else \
		echo "==> Fixing source code with gofmt..." ; \
		gofmt -s -w ./$(pkg) ; \
	fi

# Install dev lints tools 
tools:
	@echo ""
	@echo "==> Installing missing commands dependencies..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_VERSION)

# it runs the lints
lint: fmt
## this is the same that make checktools command to validate the tools
	@if ! sh -c "'./scripts/checktools.sh'" 2>&1 /dev/null ; then \
		echo "" ; \
		echo "==> There are missing tools..." ; \
		$(MAKE) -s tools ; \
	fi
	@if [ "${pkg}" = "" ] ; then \
		echo "" ; \
		echo "==> Linting all files" ; \
		bin/golangci-lint run -v ; \
	else \
		echo "==> Checking source code against linters..." ; \
		cd ./$(pkg) ; \
		./../bin/golangci-lint run ./... -v ; \
	fi


# check all the test and run the linter
check: test lint

docker-build:
# e.g: make docker-build service=backendify version=v1.0.0-test
	docker build . \
	--no-cache \
	--progress=auto \
	--build-arg SERVICE=$(service) \
	--build-arg GO_COMMANDS="$(go_commands)" \
	-t $(service):$(version)

docker-run:
# e.g: make docker-run service=backendify args=-d version=v1.0.0-test params="ru=http://localhost:9001 us=http://localhost:9002"
	@docker run -it --rm \
	--entrypoint ./backendify \
	--name backendify \
	$(args) \
	-p 9000:9000 $(service):$(version) \
	$(params)

docker-bnr: 
# e.g: make docker-bnr service=backendify args=-d version=v1.0.0-test params="ru=http://localhost:9001 us=http://localhost:9002"
	$(eval version=$(if $(version),$(version),latest))
	echo "${version}"
	$(MAKE) docker-build version=$(version) && $(MAKE) docker-run version=$(version)
	
.PHONY: test fmtcheck checktools fmt tools lint check docker-run docker-bnr docker-build
