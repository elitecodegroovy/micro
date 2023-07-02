
# Install PostgreSQL

docker run --name local-psql -v local_psql_data:/var/lib/postgresql/data -p 54320:5432 -e POSTGRES_PASSWORD=my_password -d postgres

e.g.

```
docker run --name psql15 -v /Users/superpowerai/pgsql/data:/var/lib/postgresql/data -p  54320:5432 -e POSTGRES_USER=postgres  -e POSTGRES_PASSWORD=Postgresql101010001010101 -e POSTGRES_PASSWORD=Postgresql101010001010101 -d postgres

```


docker run --name psql15 -v /opt:/var/lib/postgresql/data -p  5433:5432 -e POSTGRES_USER=postgres  -e POSTGRES_PASSWORD=Postgresql101010001010101 -d postgres

