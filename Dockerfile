FROM golang:1.15-alpine AS build
WORKDIR /src
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux 
ENV GOARCH=amd64

RUN go mod download
RUN go build -o /bin/api ./main.go
RUN go build -o /bin/migratedb cmds/migrate.go

FROM scratch AS bin
COPY --from=build /bin/api /api
COPY --from=build /bin/migratedb /migratedb

CMD ["/api"]
