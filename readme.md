# Go Bank API

First API in Golang!

## Endpoints

- GET - /account (Get All Accounts)
- GET - /account/{id} (Get Account By Id - Needs Autentication) 
- POST - /account (Create New Account)
- DELETE - /account/{id} (Delete account - Needs Autentication)
- PATCH - /transfer (Transfer Amount to especified account - Needs Autentication)
- POST - /login (Create JWT Token)

## Docs

The authentication are from an custom header:  `x-jwt-token`
The delete and detail just the owner account can request
The transfer endpoint validades account balance if the accountID are diferrent of the request accountID.  

## How to Run

1. `docker compose up -d`
