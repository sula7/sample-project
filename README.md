## This is a sample project for interviewers

An HTTP server written on Golang. Main storage is Postgres and token storage is Redis

## Env vars

`DB_HOST` IP or domain name where database locates   
`DB_NAME` Database which will use this app  
`DB_PORT` Port that database listens  
`DB_USER` Credential "user" to get access  
`DB_PASSWORD` Credential "password" to get access  
`TOKEN_SECRET` A secret for JWT token (signature)  
`REDIS_PASSWORD` Credential to access redis DB (token storage)

## Prepare to start

Need postgres from 9.6 to 13 before start app.  
This project contains `docker-compose.yml` file to execute DB:
run `docker-compose up` from project root dir

## App run flags

*NOTE:* No flags means the app will start migrating all database versions  
-db-version=[int] Specifies to which version need to migrate database
