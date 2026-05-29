# ==============================================================================
# GMAO Backend Makefile - Local Developer testing
# ==============================================================================

.PHONY: help infra infra-down run status

help:
	@echo "Available commands:"
	@echo "  make infra       - Spin up database (Postgres) and service discovery (Consul) in Docker"
	@echo "  make infra-down  - Stop and clean the Postgres and Consul containers"
	@echo "  make run         - Launch all Go services locally using the run_all.ps1 script"
	@echo "  make status      - Print list of services currently registered in Consul"

infra:
	docker compose -f deploy/docker-compose.yml up -d consul postgres

infra-down:
	docker compose -f deploy/docker-compose.yml down -v

run:
	powershell -ExecutionPolicy Bypass -File ./run_all.ps1

status:
	curl -s http://127.0.0.1:8500/v1/catalog/services | json_pp || curl -s http://127.0.0.1:8500/v1/catalog/services
