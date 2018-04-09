The GraphQL service only supports member-level data access at this stage.

Start dev server: 
```bash
$ go run cmd/webd/main.go -s graphql
``` 

GraphiQL available at http://localhost:5001/graphql


## Sample Queries##

Member Query acts as a *viewer* type, requiring a valid JWT token to access child nodes


**Member Query**
```graphql
query Member($token: String!) {
  memberUser(token: $token) {
    id
    lastName
    middleNames
    firstName
    email
    mobile
    locations {
      type
      address
      phone
    }
  }
}
``` 

Mutation Examples:

**Record a member activity**

Note that the activityId and TypeId must be present, and must be related.

```graphql
mutation Member($token: String!) {
  member(token: $token) {
    saveActivity(obj: {
      date:"2018-02-03"
      description: "Update the internal member activity data funcs"
      quantity: 4
      activityId: 22
      typeId: 18    
    }) 
    {
      id
      activity
      type
      category
      credit
      date
      description
    }
  }
}
```


## Visual Representation

Run the introspection query below and paste the output into https://apis.guru/graphql-voyager/ 


**Introspection Query**

```graphql
query IntrospectionQuery {
  __schema {
    queryType {
      name
    }
    mutationType {
      name
    }
    subscriptionType {
      name
    }
    types {
      ...FullType
    }
    directives {
      name
      description
      args {
        ...InputValue
      }
      onOperation
      onFragment
      onField
    }
  }
}

fragment FullType on __Type {
  kind
  name
  description
  fields(includeDeprecated: true) {
    name
    description
    args {
      ...InputValue
    }
    type {
      ...TypeRef
    }
    isDeprecated
    deprecationReason
  }
  inputFields {
    ...InputValue
  }
  interfaces {
    ...TypeRef
  }
  enumValues(includeDeprecated: true) {
    name
    description
    isDeprecated
    deprecationReason
  }
  possibleTypes {
    ...TypeRef
  }
}

fragment InputValue on __InputValue {
  name
  description
  type {
    ...TypeRef
  }
  defaultValue
}

fragment TypeRef on __Type {
  kind
  name
  ofType {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
      }
    }
  }
}
```




 
 