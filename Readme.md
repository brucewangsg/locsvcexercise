# Simple app to handle storing locations and users

## Database

The app uses postgresql to store data. Interfacing with db from go app is done using gorm. We use gorm mainly for its automigrate and basic create, update and find individual records.

### migrating database

```
go build -o bin/migratedb cmds/config.go cmds/migrate.go
DB_HOST={host} DB_USER={user} DB_NAME={dbname} DB_PORT={dbport} DB_PASS={yourdbpass} bin/migratedb
```

replace env variables accordingly