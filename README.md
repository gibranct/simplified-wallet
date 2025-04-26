# Backend Challenge

## Technologies 
- Golang
- Keycloak
- Postgres
- SNS
- SQS

## How to run

```shell
docker compose up --build -d
```

```shell
make run
```

OR 

```shell
make run/docker
```

## Endpoints

### Create user
```http request
POST /v1/users
Content-Type: application/json

{
  "fullName": "Pedro da Silva",
  "CPF": "11111111111",
  "email": "pedro@mail.com",
  "password": "",
  "passwordConfirmation": "",
  "type": "COMMON|MERCHANT"
}
```

### Get users
```http request
GET /v1/users
Content-Type: application/json

[{
  "fullName": "Pedro da Silva",
  "CPF": "11111111111",
  "email": "pedro@mail.com",
  "type": "COMMON|MERCHANT"
}]
```

### Make transfer
```http request
POST /v1/transfer
Content-Type: application/json

{
  "value": 100.0,
  "sender": "uuid",
  "receiver": "uuid"
}
```








