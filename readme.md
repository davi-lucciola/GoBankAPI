# Go Bank API
First API in Golang!

## Endpoints
- GET - /account (Get All Accounts)
- GET - /account/{id} (Get Account By Id - Needs Authentication)
- POST - /account (Create New Account)
- DELETE - /account/{id} (Delete Account - Needs Authentication)
- PATCH - /transfer (Transfer Amount to Specified Account - Needs Authentication)
- POST - /login (Create JWT Token)

## Docs

The authentication is via a custom header: x-jwt-token. 

Only the owner of the account can request the delete and detail endpoints. 

The transfer endpoint validates account balance if the accountID is different from the request accountID.

## How to Run
- `docker compose up -d`
