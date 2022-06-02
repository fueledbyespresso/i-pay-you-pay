[![Build Status](https://travis-ci.com/fueledbyespresso/school-supply-list.svg?branch=master)](https://travis-ci.com/fueledbyespresso/school-supply-list)

# Setup Project

## Initialize a postgres database

Create a postgres database*

Place postgres URL in the environment file

``DATABASE_URL=user:password@host:port/database``

Set a database secret for storing session data. Use a random password or string generator to create a secure,
unguessable secret.
``DATABASE_SECRET=<SECRET>``

*Project made using Postgres 12.3

## Create Google Oauth credentials

Go to https://console.developers.google.com/apis/dashboard and create a new project. Create a new Oauth 2.0 Client ID
and secret in the ``credentials`` tab Populate the projectvars.env file with your Oauth Client ID and secret

``GOOGLE_CLIENT_ID=<CLIENTID>``

``GOOGLE_CLIENT_SECRET=<SECRET>``

## Create Administrator

Login to your application once with a Goolge account to generate your user account

Execute the following
SQL `INSERT INTO user_role_bridge (user_uuid, role_id) SELECT user_id, 1 from account WHERE email='<YOUR_EMAIL'`

## Add Database Changes

Database changes are tracked in database/migrations folder and will be executed at compilation.
The current database version will be kept in the database table `schema_migrations`.
``migrate create -ext sql -dir database/migrations -seq <MIGRATION_NAME>``

## Development environment

Set the React proxy in the frontend/package.json file to any port. By default it points to port 5000.
Run a React development environment using `npm start`
Run the Go development environment using `go build main.go`

### Frontend

Running `react start` from the root directory will run the react app in development mode on port 3000. The backend will
be queried through port 5000 with the proxy for react listening to port 5000.

### Backend

Running `go run build` from root directory will run the backend on port 3000
