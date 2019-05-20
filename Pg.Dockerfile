FROM postgres

COPY ./db/postgres/create_table.sql /docker-entrypoint-initdb.d/create_table.sql