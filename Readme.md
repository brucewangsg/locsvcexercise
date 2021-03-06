# Simple app to handle storing locations and users

## Database

The app uses postgresql to store data. Interfacing with db from go app is done using gorm. We use gorm mainly for its automigrate and basic create, update and find individual records.

### migrating database

Spin up your postgresql database. Or you can use docker-compose to create your pg instance

```
docker-compose up -d db
```

If you would like to setup default db password

```
export POSTGRES_PASSWORD={yourdbpass}
docker-compose up -d db
```

Create databases

```
docker exec -ti locsvcexercise_db_1 createdb -U postgres locexercise
docker exec -ti locsvcexercise_db_1 createdb -U postgres locexercise_test
```


```
go build -o bin/migratedb cmds/migrate.go
DB_HOST={host} DB_USER={user} DB_NAME={dbname} DB_PORT={dbport} DB_PASS={yourdbpass} bin/migratedb
```

Replace env variables accordingly. If you are developing on your local environment, you can copy `sample.env` into `.env` file. Change all the configuration needed on that file.


## API Server

starting api server

```
go build -o bin/api main.go
DB_HOST={host} DB_USER={user} DB_NAME={dbname} DB_PORT={dbport} DB_PASS={yourdbpass} APP_PORT=5678 bin/api
```

### API endpoints

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
GET /auths/verify

e.g.
curl -X GET "http://localhost:5678/auths/verify" -H "Accept: application/json" -H "Authorization: Bearer {jwtToken}"
```

```
GET /locations

curl "http://localhost:5678/locations"
```


```
GET /locations/:id

curl "http://localhost:5678/locations/1"
```


```
POST /location_preference

e.g.
curl -X PUT "http://localhost:5678/location_preference" -H "Accept: application/json" -H "Authorization: Bearer {jwtToken}" -d '{"location_id":1}'
```

## Setting up server

```
docker-compose up -d --scale app=2

# create databases

docker exec -ti locsvcexercise_db_1 createdb -U postgres locexercise

# migrate schema and seed data

docker exec -ti locsvcexercise_app_1 /migratedb
```

Test it out

```
# go to http://localhost:5678/, and you will see nothing

curl -X POST "http://localhost:5678/auths/register" -d '{"name": "John Chow", "email": "john@email.com", "password": "PASS"}'
curl -X POST "http://localhost:5678/auths/login" -d '{"email": "john@email.com", "password": "PASS"}'

# get the jwt token from the response, the token will be used for auth token for subsequent requests
# or you can use jq

token=$(curl --silent -X POST "http://localhost:5678/auths/login" -d '{"email": "john@email.com", "password": "PASS"}' | jq --raw-output '.Token')

curl "http://localhost:5678/locations"
curl -X PUT "http://localhost:5678/location_preference" -H "Authorization: Bearer $token" -d '{"location_id":1}'

# you will see saved location in your profile

curl "http://localhost:5678/location_preference" -H "Authorization: Bearer $token"
```

Bonus

```
# book a location

curl -X PUT "http://localhost:5678/bookings/1" -H "Authorization: Bearer $token"
```