# Fold-ELK

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

### Setting up the environment

- Clone the repo
- Run `mv .env.example .env` to create the environment file
- Run `docker-compose up --build` to build and bootup the ELK stack
- Once the ELK stack is up and running, you can use the below 
