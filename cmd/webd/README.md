## webd/

Contains executable for the web server api.

Can start as REST server, or a GraphQL server. This is a workaround for the one web-process 
limit on Heroku and allows the same repo to be pushed to two separate Heroku apps - one for 
the REST server and the other GraphQL.



**flags**

`-s` server type, values can be 'rest' or 'graphql'

`-p` port number, defaults to 5000 for rest server, 5001 for graphql server


eg. Start graphql server on dev machine:

```bash
$  go run cmd/webd/main.go -s graphql
```

Alternatively, set env var:
 
```
WEBD_TYPE=[rest|graphql]
```

