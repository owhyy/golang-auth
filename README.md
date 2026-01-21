# golang-blog

Simple blog application featuring user auth and blog post CRUD

## Configuration

See `.env.example` for configuration options

## Running

- Locally:

```
go run ./cmd/web
```

- Dockerfile (recommended):

```
docker build -t go-blog-app .
docker run --env-file .env --rm -p 8080:8080 go-blog-app:latest
```

Accessing localhost:8080 should open the home page.

> NOTE: Keep in mind that restarting the container will reset the database and fileserver since containers don't have persistent storage.

## Useful commands

There are some scripts in the cmd/debug folder that could prove useful:

- Generating some fake data can be done by using the `populate` command. By default it will generate 1000 posts and 10 users, but this can be configured by passing specific flags. To simplify testing, it will also generate an admin account with credentials admin@example.com and password admin.

```
go run ./cmd/debug/populate.go
```

- Creating a admin account is done via the `createadmin` command. The credentials are passed via command line arguments, like this:

```
go run ./cmd/debug/createadmin.go -email <email> -password <password> -username <username>
```

### Running under docker

To run the commands under docker, you need to have a running container. Binaries of `populate` and `createadmin` are available in the root of the container and can be run like so.

```
docker exec <container-name> populate
```

Respectively for `createadmin`:

```
docker exec <container-name> createadmin -email <email> -password <password> -username <username>
```