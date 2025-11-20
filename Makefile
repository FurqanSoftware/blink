# Copyright 2025 Furqan Software Ltd. All rights reserved.

.PHONY: lint
lint:
	staticcheck ./...

.PHONY: lint.tools.install
lint.tools.install:
	go install honnef.co/go/tools/cmd/staticcheck@2025.1.1

.PHONY: test
test:
	go test -v ./...
