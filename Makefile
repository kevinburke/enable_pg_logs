RELEASE := $(GOPATH)/bin/github-release

$(RELEASE):
	go get -u github.com/aktau/github-release

release: $(RELEASE)
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	git tag $(version)
	git push origin --tags
	mkdir -p releases/$(version)
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/enable_pg_logs-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/enable_pg_logs-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/enable_pg_logs-windows-amd64 .
	# these commands are not idempotent so ignore failures if an upload repeats
	$(RELEASE) release --user kevinburke --repo enable_pg_logs --tag $(version) || true
	$(RELEASE) upload --user kevinburke --repo enable_pg_logs --tag $(version) --name enable_pg_logs-linux-amd64 --file releases/$(version)/enable_pg_logs-linux-amd64 || true
	$(RELEASE) upload --user kevinburke --repo enable_pg_logs --tag $(version) --name enable_pg_logs-darwin-amd64 --file releases/$(version)/enable_pg_logs-darwin-amd64 || true
	$(RELEASE) upload --user kevinburke --repo enable_pg_logs --tag $(version) --name enable_pg_logs-windows-amd64 --file releases/$(version)/enable_pg_logs-windows-amd64 || true
