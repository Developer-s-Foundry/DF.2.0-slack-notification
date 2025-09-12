.SILENT:

POSTGRES_CONTAINER=my_postgres
POSTGRES_USER=admin
POSTGRES_PASSWORD=secret
POSTGRES_DB=mydb
POSTGRES_PORT=5432
POSTGRES_IMAGE=postgres:15
BINARY_NAME=app

.PHONY: postgres start stop logs psql remove

postgres:
	docker run -d \
		--name $(POSTGRES_CONTAINER) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 \
		$(POSTGRES_IMAGE)

start:
	docker start $(POSTGRES_CONTAINER)

stop:
	docker stop $(POSTGRES_CONTAINER)

remove:
	docker rm -f $(POSTGRES_CONTAINER)

logs:
	docker logs -f $(POSTGRES_CONTAINER)

psql:
	docker exec -it $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

build: 
	go build -o $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)