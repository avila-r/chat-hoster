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
