.PHONY: help# Generate the listo of targets with description
help:
	@echo "Available make targets:"
	@awk	'/^.PHONY: .* #/ { \
			taregt = substr($$0, 10, index($$0, " #") - 10); \
			helpText = substr($$0, index($$0,"# ") + 2); \
			printf "%-20s %s\n", target,helpText; \
		}' Makefile

.PHONY: build # Build a specific binary
build:
	@echo "What binary do you want to build?"
	@read -p "name: " NAME \
	&& go build -o cmd/$${NAME}/bin/orch-$${NAME} -v ./cmd/$${NAME}/
	@echo "Done."

.PHONY: build-all # Build all specific binary
build-all:
	@echo "Building all binaries..."
	@go build -o cmd/gateway/bin/orch-gateway -v ./cmd/gateway/
	@echo "Done."

.PHONY: swag # Generate Swagger files
	@echo "Generate swagger files..."
	@swag init -g swagger.go --parseDependency

.PHONY: testrun # Run test and generate coverage
testrun:
	@echo "running tests.."
	@go test ./... -coverprofile=./coverage.out -coverpkg ./...

.PHONY: testrunv # run tests verbosely and generate coverage
testrunv:
	@echo "running tests..."
	@go test -v ./... -coverprofile=./coverage.out -coverpkg

.PHONY: cleandbs # remove all .db files generated from tests
cleandbs:
	@find . -type f -name "*.db" -delete
	@echo "Deleted all .db files"

.PHONY: test # run tests and generate coverage after cleaning
test: testrun cleandbs

.PHONY: testv # run tests verbosely and generate coverage after cleaning
testv: testrunv cleandbs

.PHONY: report # show test coverage reports
report:
	@echo "showing test coverage report....."
	@go tool cover -html=coverage.out

.PHONY:	db #	Build Docker compose images.
db:
	@docker compose up -d db_for_app

