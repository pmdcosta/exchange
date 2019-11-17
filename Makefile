.PHONY: \
		build \
		test \
		start \
		stop

SHELL := /bin/bash

# Run Unit tests
test:
	docker-compose -f docker-compose.test.yml up --build

# Builds the Docker container
build:
	docker build -t pmdcosta/exchange:latest .

# Starts an instance of the exchange
start:
	docker-compose up -d

# Stops the exchange instance
stop:
	docker-compose down