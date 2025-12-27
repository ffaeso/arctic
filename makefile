# standard test with coverage
.PHONY: test
test:
	go test ./... -v -cover

# for seeing html visualisation of whats covered and whats not
.PHONY: test-html
test-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out
