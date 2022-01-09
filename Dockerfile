FROM golang:1.17 AS build

ADD . /app
WORKDIR /app
RUN go build ./cmd/main/main.go

FROM ubuntu:20.04

RUN apt-get -y update &&\
    apt-get install -y tzdata

ENV POSTGRES_VERSION 12

RUN apt-get -y update &&\
    apt-get install -y postgresql-${POSTGRES_VERSION}

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER ekasy WITH SUPERUSER PASSWORD 'ekasy';" &&\
    createdb -O postgres forum &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/${POSTGRES_VERSION}/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/${POSTGRES_VERSION}/main/postgresql.conf

EXPOSE 5432

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .

COPY --from=build /app/main .

EXPOSE 5000

ENV PGPASSWORD ekasy

CMD service postgresql start &&\
    psql -h localhost -d forum -U ekasy -p 5432 -a -q -f ./db/init.sql &&\
    ./main
