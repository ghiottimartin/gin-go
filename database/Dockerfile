FROM postgres:10.3

# Create volume to store data

COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]