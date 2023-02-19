ARG ELASTIC_VERSION

FROM docker.elastic.co/logstash/logstash:${ELASTIC_VERSION}

RUN /opt/logstash/bin/logstash-plugin install logstash-integration-jdbc
RUN mkdir /opt/logstash/vendor/jdbc
RUN curl -o /opt/logstash/vendor/jdbc/postgresql-42.2.18.jar https://jdbc.postgresql.org/download/postgresql-42.2.18.jar

ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]
