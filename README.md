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

## Running the tests

To run the tests, you can use the following command:

```bash
make test
```