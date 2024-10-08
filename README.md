## Book keeping REST API in Go 

### Description
 library API project is a Go-based application designed to manage a digital library. It supports book creation, updates, deletion, and retrieval with JWT authentication for both users and admins. The API allows users to search books, post reviews, view reviews, and add books to their favorites list. Admins can manage the library by adding or removing books and reviews. The project also integrates Elasticsearch for enhanced book search functionality and uses MySQL as the primary database, with routes built using the Gorilla Mux router.

### Installation

Open a PowerShell terminal (version 5.1 or later) and from the PS C:\> prompt, run: 

```bash 
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression

```
Once installation is done install 

```bash
scoop install main/make
```

and 

```bash
scoop install main/gcc
```

Relaunch all terminals 

Also make sure you have the following tools installed on your machine.

- [Migrate (for DB migrations)](https://github.com/golang-migrate/migrate/tree/v4.17.0/cmd/migrate)

## Running the project

Firstly make sure you have a MySQL database running on your machine or just swap for any storage you like under `/db`.

Then create a database with the name you want *(`ecom` is the default)* and run the migrations. Before you do so,

Locate the `.env.example` file and copy it content inside a new file `.env`

```env
# Server
PUBLIC_HOST=http://localhost
PORT=8080

# Database
DB_USER=root
DB_PASSWORD=mypassword
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=ecom
```

Change the password `mypassword` to your MYSQL localhost password.


Create the Database Manually:

    Open the MySQL command line:

```bash
mysql -u root -p
```

After logging in, create the ecom database:

```sql

CREATE DATABASE ecom;
```
Verify that the database was created:


```sql
SHOW DATABASES;
```

You should see ecom listed in the output.

2. Run the Migration:

    Now that the ecom database exists, try running the make migrate-up command:

```bash
make migrate-up
```

After that, you can run the project with the following command:

```bash
make run
```

## Running the tests

To run the tests, you can use the following command:

```bash
make test
```


### Running the Project with Docker
#### Step 1: Docker Setup

To use Docker for containerizing the application and the database, follow these steps:

- Ensure you have Docker installed on your machine.

- Create a docker-compose.yml file with the following content:

```bash

version: '3'
services:

  db:
    image: mysql:8.0
    healthcheck:
      test: "exit 0"
    volumes:
      - db_data:/var/lib/mysql
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: mypassword
      MYSQL_DATABASE: ecom

  api:
    build: 
      context: .
      dockerfile: Dockerfile
    restart: on-failure
    volumes:
      - .:/go/src/api
    ports:
      - "8080:8080"
    environment:
      DB_HOST: db
      DB_USER: root
      DB_PASSWORD: mypassword
      DB_NAME: ecom
    links:
      - db
    depends_on:
      - db

volumes:
  db_data:
```
Don't forget to change `mypassword` to your database password

Create a Dockerfile in your project root:


```bash

    # syntax=docker/dockerfile:1

    # Build the application from source
    FROM golang:1.22.0 AS build-stage
      WORKDIR /app

      COPY go.mod go.sum ./
      RUN go mod download

      COPY . .

      RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

      # Run the tests in the container
    FROM build-stage AS run-test-stage
      RUN go test -v ./...

    # Deploy the application binary into a lean image
    FROM scratch AS build-release-stage
      WORKDIR /

      COPY --from=build-stage /api /api

      EXPOSE 8080

      ENTRYPOINT ["/api"]
``` 

#### Step 2: Running the Docker Containers

To start the application and database containers, use:

```bash
docker-compose up --build
```

This command will build and start both the MySQL database and the API service. The application will be accessible at http://localhost:8080.



## API Endpoints
User Registration
- Endpoint: POST /api/v1/register
- Description: Registers a new user.
- Payload Example:

```bash
{
  "email": "jdmasciano2@gmail.com",
  "password": "ogbono",
  "firstName": "ade",
  "lastName": "burger"
}
```

Admin Registration: Include "role": "admin" in the payload.

```bash

    {
      "email": "jd@gmail.com",
      "password": "ogbono",
      "firstName": "surfinia",
      "lastName": "burg",
      "role": "admin"
    }
```


User Login

- Endpoint: POST /api/v1/login

- Description: Logs in a user and returns a JWT token.

- Payload Example:


```bash
{
  "email": "me@me.com",
  "password": "asd"
}
```

Response Example:

```bash
json

    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOjE3MjU5MTk0NTYsInVzZXJJRCI6IjEifQ.Ww85HQzCdhzp_LzJTg8UvxcrXj5eanLyLJDyDNQIG6E"
    }
```


Get User by ID

- Endpoint: GET /api/v1/users/{id}

- Description: Retrieves a user's information by their ID.

- Response Example:

```bash

    {
      "id": 4,
      "firstName": "ade",
      "lastName": "surfinia",
      "email": "me@ade.com",
      "role": "user",
      "createdAt": "2024-09-01T20:12:23Z"
    }
``` 

You can find more details on the Library API endpoints here (`https://docs.google.com/document/d/1vZ_6MTWN3PK9Ol02pT37M_--D5dRimy-XgiPwgoDlkM/edit?usp=sharing`)