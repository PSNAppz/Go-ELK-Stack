# Fold-ELK

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/12298080-b2a522f0-5323-4aa1-b803-c73160cde976?action=collection%2Ffork&collection-url=entityId%3D12298080-b2a522f0-5323-4aa1-b803-c73160cde976%26entityType%3Dcollection%26workspaceId%3D42654ab9-e148-4e67-b4b3-90fd805dfb7f)

The postman collection contains all the API endpoints along with examples.
### The objectives of this repo are as follows:

- To build a data pipeline that syncs and transforms the data stored in a PostgresSQL Database to Elasticsearch, making the data searchable in a fuzzy and relational manner.
- To build a REST API service that will query Elasticsearch.

### Project brief

Fold-ELK has a set of configuration files for the ELK stack (Elasticsearch, Logstash, Kibana). You can follow the below steps to get the ELK stack up and running.
CRUD APIs are built using go-gin and the data is stored in PostgresSQL. The data is then synced to Elasticsearch using Logstash. The APIs are then used to query Elasticsearch.
Data is seeded using faker, which is used to populate the PostgresSQL database.

## Getting Started

### Prerequisites

- Docker
- Docker Compose
- Go

### Setting up and building the docker environment

- Clone the repo
- Run `mv .env.example .env` to create the environment file
- Run `docker-compose up --build` to build and bootup the ELK stack
- Once the ELK stack is up and running, you can use the below 
- Kibana: http://localhost:5601
- Elasticsearch: http://localhost:9200
- API: http://localhost:8080

### Running the API locally

The CRUD API can be run locally using the below steps:

- Make sure the ELK stack is up and running using the above steps
- Run `go build -o bin/application cmd/api/main.go` to build the API binary
- Run `export PGURL="postgres://fold-elk:password@localhost:5432/fold_elk?sslmode=disable"`
- Migrate the database using `migrate -database $PGURL -path db/migrations/ up `
- Run `./bin/application` to run the API and start listening on port 8080

### Pre Submission Checklist
- [x] Detailed steps are included that allow us to spin up the services and the data pipeline.
- [x] [Loom video URL](https://www.loom.com/share/9b76a3cf38cf4a48b40936adae8e74e9)
