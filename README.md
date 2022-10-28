[![LinkedIn Badge](https://img.shields.io/badge/LinkedIn-Profile-informational?style=flat&logo=linkedin&logoColor=white&color=0D76A8)](https://www.linkedin.com/in/daniil-drozdov-a5393521b/)

# Виконання тестового завдання Trainee Golang Developer в EVO

## Practice project

> **Note**
> Dependencies

* Docker
* docker-compose
* Make

### Quick start

> **Note**
> First of all create and configure ops/.env file with next variables

```
CSVAPP_DBURI="DB_URI"
CSVAPP_ADDR=":APPLICATION_PORT" Default: ":8080"
CSVAPP_DEBUG=TRUE_OR_FALSE

POSTGRES_USER=DB_USER
POSTGRES_PASSWORD=DB_PASSWORD
POSTGRES_DB=DB_NAME
```

> **Warning**
> If CSVAPP_ADDR is not default don't forget to change it to new port in docker-compose.yml "4444:`new_port`"

---

**Run**
```shell
make
```

> **Note**
> By default application will be hosted on http://localhost:4444

> **Note**
> By default documentation will be hosted on http://localhost:4444/swagger/index.html

---

> **Warning**
> For **generate**, **generate-sql** and **generate-swagger** you use next command

Installation base package
```shell
make init
```

Remove containers
```shell
make rm
```

---

Generate SQLC code and swagger documentation
```shell
make generate
```

### Useful make command

Show logs application container
```shell
make logs
```

> **Note**
> If you change db/migrations/ or queries/ for regenerate SQLC code

Generate SQLC 
```shell
make generate-sql
```

> **Note**
> If you change documentations for regenerate swagger documentation

Generate swagger documentation
```shell
make generate-swagger
```

Format swagger documentation in code
```shell
make fmt
```