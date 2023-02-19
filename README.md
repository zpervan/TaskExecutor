# TaskExecutor #

A simple application which executes tasks from the queue received via a CLI request.

The application consists of the following components:
* Agent         - Executes the received commands
* API (Handler) - Receive and process the received tasks (requests)
* Database      - A Mongo database which stores all received tasks

Technologies used:
* Golang
* Docker

## Setup ##

In order to run this application successfully, make sure you have installed [docker](https://docs.docker.com/get-docker/) on your PC.
Make sure that [curl](https://curl.se/) is installed on your PC in order to send request to the application.

## Usage ##

### Production ###

Position your terminal into the root of the project and execute the following command:
```shell
$ docker-compose up -d
```

* Add a new task:
```shell
$ curl --request POST --data '{"command": "echo Hello && sleep 10"}' -H 'Content-Type: application/json' http://localhost:3500/tasks
```

* Get all tasks:
```shell
$ curl localhost:3500/tasks
```

### Development or testing ###

In order to develop or test the functionality locally, you can run each package separately. 

Example: Testing the backend API

* Position your terminal at the root of the project
* You can use the database installed on your local PC or even run the `database` Docker image standalone by executing:
```shell
$ docker-compose up -d --build database
```
* In order to successfully run the backend API localhost, make sure that we use the correct database URL, i.e. `mongodb://admin:password@localhost:27018`
* Run the backend API by executing:
```shell
$ go run taskexecutor/backend
```
* You can now test the backend by sending commands as previously explained

NOTE: For more details about the exposed ports or network in general, have a look into the `docker-compose.yaml` file
