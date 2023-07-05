# **UserManagementBackend**

## **How to run / execute the project**

```bash
# you need to clone the project
$ git clone https://github.com/emil-1003/UserManagementBackendGolang.git

# create .env file in root with the following variables:
DB_NAME
DB_USERNAME
DB_PASSWORD
DB_HOST
DB_PORT
TOKEN_SECRET_WORD

# run project
$ go run main.go
```

## **How to setup project in docker**
```bash
# Build image
$ cd project/Dockerfile

$ sudo docker build -t [name] .

# Start container by running image
# 1. --name set name
# 2. -p set ports
# 3. -d run in background
$ docker run --name=[name] -p 80:8888 -d [image name]
```

---
> *Created by - Emil Andersen - 2023*
