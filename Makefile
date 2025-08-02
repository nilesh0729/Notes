Container:
	docker run --name ganu -p 5431:5432 -e POSTGRES_PASSWORD=MummyJi -e POSTGRES_USER=root -d postgres

CreateDB:
	docker exec -it ganu createdb --username=root --owner=root PapaJi

DropDB:
	docker exec -it ganu dropdb -U root PapaJi

MigrateUp:
	migrate -path db/migrate -database "postgres://root:MummyJi@localhost:5431/PapaJi?sslmode=disable" -verbose up

MigrateDown:
	migrate -path db/migrate -database "postgres://root:MummyJi@localhost:5431/PapaJi?sslmode=disable" -verbose down

Sqlc:
	sqlc generate

.PHONY:	Container	CreateDB	DropDB	MigrateUp	MigrateDown	Sqlc