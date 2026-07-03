# ----------------------------------------------------------MIGRATIONS----------------------------------------------------------------------------

migrate:
	bash deployments/scripts/migration.sh

migrateup:
	migrate -path deployments/migrations/up -database "mysql://root:password@tcp(localhost:3306)/falcon" up

migratedown:
	migrate -path deployments/migrations/down -database "mysql://root:password@tcp(localhost:3306)/falcon" down

# ----------------------------------------------------------RUN SERVICES SCRIPT-------------------------------------------------------------------

run-core:
	cd ./services/core/app && wire
	cd ./services/core/app/ && go run .

# ----------------------------------------------------------GENERATE PROTOC FILE------------------------------------------------------------------
