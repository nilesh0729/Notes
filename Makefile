DB_URL=postgres://root:MummyJi@127.0.0.1:5431/PapaJi?sslmode=disable

Container:
	docker run --name ganu -p 5431:5432 -e POSTGRES_PASSWORD=MummyJi -e POSTGRES_USER=root -d postgres

Createdb:
	docker exec -it ganu createdb --username=root --owner=root PapaJi

Dropdb:
	docker exec -it ganu dropdb -U root PapaJi

MigrateUp:
	migrate -path internal/db/migrate_files -database "$(DB_URL)" -verbose up

MigrateDown:
	migrate -path internal/db/migrate_files -database "$(DB_URL)" -verbose down

Sqlc:
	sqlc generate

Test:
	go test -v -cover ./...

Mock:
	mockgen -package mockDB -destination internal/db/Mock/gomock.go github.com/nilesh0729/Notes/internal/db/Result Store

Server:
	go run cmd/api/main.go

.PHONY:	Container	Createdb	Dropdb	MigrateDown	MigrateUp	Sqlc	Server