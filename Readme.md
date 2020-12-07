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
docker exec -ti locsvcexercise_db_1 createdb -U postgres locexercise_test
```

```
go build -o bin/migratedb cmds/migrate.go
DB_HOST={host} DB_USER={user} DB_NAME={dbname} DB_PORT={dbport} DB_PASS={yourdbpass} bin/migratedb
```

replace env variables accordingly


## API Server

starting api server

```
go build -o bin/api main.go
DB_HOST={host} DB_USER={user} DB_NAME={dbname} DB_PORT={dbport} DB_PASS={yourdbpass} APP_PORT=5678 bin/api
```

### endpoints

```
POST /auths/register

e.g.
curl -X POST "http://localhost:5678/auths/register" -H "Accept: application/json" -d '{"email": "user@gmail.com", "name": "User 1", "password": "PASS"}'
```


```
POST /auths/login

e.g.
curl -X POST "http://localhost:5678/auths/login" -H "Accept: application/json" -d '{"email": "user@gmail.com", "password": "PASS"}'
```


```
POST /auths/verify

e.g.
curl -X POST "http://localhost:5678/auths/verify" -H "Accept: application/json" -H "Authorization: Bearer {jwtToken}"
```

```
GET /locations

curl "http://localhost:5678/locations"
```


```
GET /locations/:id

curl "http://localhost:5678/locations/1"
```