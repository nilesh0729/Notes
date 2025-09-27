Container:
	docker run --name ganu -p 5431:5432 -e POSTGRES_PASSWORD=MummyJi -e POSTGRES_USER=root -d postgres

Createdb:
	docker exec -it ganu createdb --username=root --owner=root PapaJi

Dropdb:
	docker exec -it ganu dropdb -U root PapaJi

MigrateUp:
	migrate -path db/migrate_files -database "postgres://root:MummyJi@127.0.0.1:5431/PapaJi?sslmode=disable" -verbose up

MigrateDown:
	migrate -path db/migrate_files -database "postgres://root:MummyJi@127.0.0.1:5431/PapaJi?sslmode=disable" -verbose down

Sqlc:
	sqlc generate

.PHONY:	Container	Createdb	Dropdb	MigrateDown	MigrateUp	Sqlc