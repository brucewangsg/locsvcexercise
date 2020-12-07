# Simple app to handle storing locations and users

## Database

The app uses postgresql to store data. Interfacing with db from go app is done using gorm. We use gorm mainly for its automigrate and basic create, update and find individual records.

### migrating database

Spin up your postgresql database. Or you can use docker-compose to create your pg instance

```
docker-compose up -d db
```

Create database schema

```
docker exec -ti locsvcexercise_db_1 createdb -U postgres locexercise
```

```
go build -o bin/migratedb cmds/migrate.go
DB_HOST={host} DB_USER={user} DB_NAME={dbname} DB_PORT={dbport} DB_PASS={yourdbpass} bin/migratedb
```

replace env variables accordingly