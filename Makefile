.PHONY: help build run build-run test logs down

help:
	@echo "commands :3"
	@echo "  make build     - build all service docker files"
	@echo "  make run       - run all services without rebuild"
	@echo "  make build-run - build all services and run"
	@echo "  make test      - run tests"
	@echo "  make logs      - open docker compose logs"
	@echo "  make down      - turn off all services"

build:
	@echo "Building all services docker files :3"
	@docker compose build
	@echo "Done! ^w^"

run: 
	@echo "running all services! >w<"
	@docker compose up -d
	@echo "to check logs use 'make logs'"

build-run: 
	@echo "rebuilding and running all services! >w<"
	@docker compose up --build -d
	@echo "to check logs use 'make logs'"

test:
	@echo "running all tests... ~w~" 
	@go test ./tests/... -v

logs:
	@docker compose logs

down:
	@echo "stopping all services ^^"
	@docker compose down
	@echo "done mrreeow :3"
