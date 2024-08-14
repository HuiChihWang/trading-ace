#!/bin/bash

# Set environment variables
export APP_ENV=test
export CONFIG_FOLDER=$(pwd)/config

# migrate the database
go run migrations/main.go

# Run tests with coverage
go test -coverprofile=coverage.out -v ./src/...

# Display coverage report
coverage=$(go tool cover -func=coverage.out | grep total: | awk '{print $3}' | sed 's/%//')
echo "Coverage: $coverage%"
if [ "$(echo "$coverage < 60" | bc)" -eq 1 ]; then
  echo "Coverage is below 60%. Failing the build."
  exit 1
fi