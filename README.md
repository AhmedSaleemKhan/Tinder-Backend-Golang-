# Simple Bank

The application that we’re going to build is a simple bank. It will provide APIs for the frontend to do following things:

1. Create and manage bank accounts, which are composed of owner’s name, balance, and currency.
2. Record all balance changes to each of the account. So every time some money is added to or subtracted from the account, an account entry record will be created.
3. Perform a money transfer between 2 accounts. This should happen within a transaction, so that either both accounts’ balance are updated successfully or none of them are.

## Setup local development

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    brew install golang-migrate
    ```

- [DB Docs](https://dbdocs.io/docs)

    ```bash
    npm install -g dbdocs
    dbdocs login
    ```

- [DBML CLI](https://www.dbml.org/cli/#installation)

    ```bash
    npm install -g @dbml/cli
    dbml2sql --version
    ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

### Setup infrastructure

- Create the bank-network

    ``` bash
    make network
    ```

- Start postgres container:

    ```bash
    make postgres
    ```

- Create simple_bank database:

    ```bash
    make createdb
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```

- Run db migration up 1 version:

    ```bash
    make migrateup1
    ```

- Run db migration down all versions:

    ```bash
    make migratedown
    ```

- Run db migration down 1 version:

    ```bash
    make migratedown1
    ```

### Documentation

- Generate DB documentation:

    ```bash
    make db_docs
    ```

- Access the DB documentation at [this address](https://dbdocs.io/techschool.guru/simple_bank). Password: `secret`

### How to generate code

- Generate schema SQL file with DBML:

    ```bash
    make db_schema
    ```

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

- Generate DB mock with gomock:

    ```bash
    make mock
    ```

- Create a new db migration:

    ```bash
    migrate create -ext sql -dir db/migration -seq <migration_name>
    ```

### How to run

- Run server:

    ```bash
    make server
    ```

- Run test:

    ```bash
    make test
    ```

## Deployment
Simple bank application is deployed on [Linode](https://www.linode.com/).

### Deploy for Staging
Simple bank application deployment for `staging` is [here](https://cloud.linode.com/linodes/40933549).

Deployed link for staging is:
```bash
http://<staging-server-ip>:9090
```

### Deploy for Production
Simple bank application deployment for `production` is [here](https://cloud.linode.com/linodes/40943446).

Deployed link for production is:
```bash
http://<prod-server-ip>:9090
```

## How to Deploy
First, We open [Linode](https://login.linode.com) website and login with the `Username` and `Password` of the Linode account. After login create a new Linode if not exist. Next, open the Linode in which you want to deploy your APIs. After opening Linode copy the `SSH Access` and run in our terminal like this:
```bash
$ ssh root@<server-ip>
```
After running this command this will require a secret `root password` of that Linode server and then entered password into the terminal of that Linode.

The next step is to install these in Linode terminal:
- [Docker Engine](https://docs.docker.com/engine/install/ubuntu/)
- [Docker Compose](https://docs.docker.com/compose/)

Next we clone this [Repository](https://github.com/flingzydev/flingzy-be)(`flingzy-be`) in Linode terminal by following this commands:
```bash
$ git clone https://github.com/flingzydev/flingzy-be.git
```

After cloning, Enter into project `flingzy-be` by this command:
```bash
$ cd flingzy-be
```
Next, We set the database certificate with that command:
```bash
$ cp <certificate.crt> /etc/ssl/certs
```
Next, set the environment for `staging` or `production` like that:
```bash
$ export ENV=env.<staging/production>
``` 

Next we deploy your APIs with `Docker Compose` command:
```bash
$ docker compose up -d
```

## Update deployment
If your wants to deploys new develop APIs or updated APIs by following this steps.

First, We open [Linode](https://login.linode.com) website and login with the `Username` and `Password` of the Linode account. After login open Linode in which you want to update and copy the `SSH Access` and run in our terminal like this:
```bash
$ ssh root@<linode-ip-address>
```
After running this command this will require a password of that Linode and then entered into the terminal of that Linode.

Next Enter into project `flingzy-be` by this command:
```bash
$ cd flingzy-be
```
And then pull new changing in `flingzy-be` with this command:
```bash
$ git pull origin master
```
If database certificate not set then add certificate with this command:
```bash
$ cp <certificate.crt> /etc/ssl/certs
```
Next, set the environment for `staging` or `production` like that:
```bash
$ export ENV=env.<staging/production>
``` 
In the last deploy your updated APIs with `Docker Compose` command:
```bash
$ docker compose up --build -d
``` 


