# env-cracker-web

This app let you upload proprietary file with `.env` extension and download the zip with all the content

![Image of the app](./app.png)

## Run project with docker

1. run :  `docker pull nicolasbellanich/env-cracker-web:latest`
1. run :  `docker run -p 8080:8080 nicolasbellanich/env-cracker-web:latest`


## Run project locally

1. Clone this repo
2. run : `go run cmd/api/main.go`