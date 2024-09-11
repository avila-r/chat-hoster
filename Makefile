container_name = go-chat-db
database_name = go-chat
database_user = root
database_password = 123

run:
	go run cmd/main.go

db-init:
	docker run --name $(container_name) -p 5432:5432 -e POSTGRES_USER=$(database_user) -e POSTGRES_PASSWORD=$(database_password) -d postgres:16.1

db-enter:
	docker exec -it $(container_name) psql

db-create:
	docker exec -it $(container_name) createdb --username=root --owner=$(database_user) $(database_name)

db-drop:
	docker exec -it $(container_name) dropdb go-chat

db-delete:
	docker stop $(container_name) || true
	docker rm $(container_name) || true

db-start:
	@if [ "`docker ps -aq -f name=$(container_name)`" ]; then \
		echo "Iniciando o container $(container_name)..."; \
		docker start $(container_name); \
		echo "Container $(container_name) iniciado com sucesso."; \
	else \
		echo "Container $(container_name) não existe. Por favor, crie o container primeiro."; \
	fi

db-stop:
	@if [ "`docker ps -q -f name=$(container_name)`" ]; then \
		echo "Parando o container $(container_name)..."; \
		docker stop $(container_name); \
		echo "Container $(container_name) parado com sucesso."; \
	else \
		echo "O container $(container_name) já está parado ou não existe."; \
	fi
