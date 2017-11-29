## webd/

Contains executable for the web server api.

If env var `GRAPHQL_SERVER=true` then it will start the GraphQL server, if not, it will start the REST server.

This is a workaround for the one web process limit on Heroku, allowing for the same repo to be pushed to two separate Heroku apps.
