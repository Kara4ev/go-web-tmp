.PHONY: keypair migrate-create migrate-up migrate-down migrate-force

PWD = $(shell pwd)
BACKPATH = $(PWD)/backend
MPATH = $(BACKPATH)/migration
PORT = 5432

N = 1

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(BACKPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(BACKPATH)/rsa_private_$(ENV).pem -pubout -out $(BACKPATH)/rsa_public_$(ENV).pem

migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(PWD)/$(APPPATH)/migrations  -seq -digits 5 $(NAME)

migrate-up:
	migrate -source file://$(APPPATH)/migrations -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable up $(N)

migrate-down:
	migrate -source file://$(APPPATH)/migrations -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable down $(N)

migrate-force:
	migrate -source file://$(APPPATH)/migrations -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable force $(VERSION)

init:
	docker-compouse up -d postgressql && \
	$(MAKE) create-keypair ENV=dev && \
	$(MAKE) create-keypair ENV=test && \
	$(MAKE) migrate-down APPPATH=backend && \
	$(MAKE) migrate-up APPPATH=backend && \
	docker-compouse down