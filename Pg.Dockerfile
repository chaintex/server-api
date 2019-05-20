FROM postgres

COPY ./config/create_table.sql /docker-entrypoint-initdb.d/create_table.sql