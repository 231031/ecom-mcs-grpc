FROM postgres:16.6-alpine

COPY up.sql /docker-entrypoint-initdb.d/

CMD ["postgres"]