## E-commerce REST API in Go 

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

List All Products

- Endpoint: GET /api/v1/products

- Description: Retrieves a list of all available products.

- Response Example:

```bash

    [
      {
        "id": 1,
        "name": "Wireless Headphones",
        "description": "High-quality wireless headphones with noise cancellation.",
        "image": "https://example.com/images/wireless-headphones.jpg",
        "price": 59.99,
        "quantity": 96,
        "createdAt": "2024-09-01T22:04:50Z"
      },
      {
        "id": 2,
        "name": "Bluetooth Speaker",
        "description": "Portable Bluetooth speaker with 360-degree sound and waterproof design.",
        "image": "https://example.com/images/bluetooth-speaker.jpg",
        "price": 79.99,
        "quantity": 44,
        "createdAt": "2024-09-01T22:09:06Z"
      }
    ]
```

Get Product by ID

- Endpoint: GET /api/v1/products/{id}

- Description: Retrieves details of a specific product by its ID.

- Response Example:

```bash

    {
      "id": 1,
      "name": "Wireless Headphones",
      "description": "High-quality wireless headphones with noise cancellation.",
      "image": "https://example.com/images/wireless-headphones.jpg",
      "price": 59.99,
      "quantity": 100,
      "createdAt": "2024-09-01T22:04:50Z"
    }
```

Checkout Cart

- Endpoint: POST /api/v1/cart/checkout

- Description: Checks out the items in the user's cart. Users need to be authenticated to use this endpoint.

- Payload Example:

```bash

{
  "items": [
    {
      "productID": 1,
      "quantity": 2
    },
    {
      "productID": 2,
      "quantity": 3
    }
  ]
}
```

Response Example:

```bash

{
  "order_id": 3,
  "total_price": 359.95
}
``` 

## Running the tests

To run the tests, you can use the following command:

```bash
make test
```