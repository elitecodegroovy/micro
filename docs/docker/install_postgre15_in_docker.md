
## postgreSQL

https://hub.docker.com/_/postgres

## install cmd
``` 
docker pull postgres:15
```


``` 
docker run --name postgresql --restart=always -e POSTGRES_USER=postgreSQL -e POSTGRES_PASSWORD=postgreSQL20230100100010101001101 -p 15432:5432 -v /home/pg:/var/lib/postgresql/data -d postgres:15

```
In the command given above,

PostgreSQL is the name of the Docker Container.
-e POSTGRES_USER is the parameter that sets a unique username to the Postgres database.
-e POSTGRES_PASSWORD is the parameter that allows you to set the password of the Postgres database.
-p 5432:5432 is the parameter that establishes a connection between the Host Port and Docker Container Port. In this case, both ports are given as 5432, which indicates requests sent to the Host Ports will automatically redirect to the Docker Container Port. In addition, 5432 is also the same port where PostgreSQL will be accepting requests from the client.
-v is the parameter that synchronizes the Postgres data with the local folder. This ensures that Postgres data will be safely present within the Home Directory even if the Docker Container is terminated.
-d is the parameter that runs the Docker Container in the detached mode, i.e., in the background. If you accidentally close or terminate the Command Prompt, the Docker Container will still run in the background.
Postgres is the name of the Docker image that was previously downloaded to run the Docker Container.


## create a database

``` 
create database pgdb
    with owner "postgreSQL" tablespace pg_default;

comment on database pgdb is 'pg15';

```