The GraphQL service only supports member-level data access at this stage.

Start dev server: 
```bash
$ go run cmd/webd/main.go -s graphql
``` 

GraphiQL available at http://localhost:5001/graphql


### Naming Conventions

The GraphQL library is quite verbose and requires a lot of related
components working together. Thus there there is a lot of opportunity
for name collisions and confusion.

In order to make the code more readable, and manaegable, the following
conventions have been used.

**File names**
* `schema` folder contains separate package folders for each of the main data entities, eg
* The `.go` file of the same name contains local types and data manipulation functions, eg `member.go`, 
and for simple types contains 
* For simple types 
*

```bash
schema/
├── activity/                 # simple type in one file
│   └── activity.go           # contains types, data funcs, graphql fields
│    
└── member/                   # complex type in separate files
    ├── activity.go           # sub-types and data functions
    ├── activitymutation.go   # mutation fields for sub-type
    ├── activityquery.go      # query fields for sub-type
    ├── member.go             # types and data functions
    ├── membermutation.go     # starting point for mutation fields
    └── memberquery.go        # starting point for query fields

```

**Variable names**

Where practical names correspond to the query graph. 

For example, the following query:

```GraphQL

query Member($token: String!) {
  member(token: $token) {
    evaluations {
      creditObtained
    }
  }
}

```

...corresponds to a field var named `queryMemberEvaluationsField` which 
returns a list of `queryMemberEvaluationObject`, as shown in this abbreviated
code example:

```Go

// field / resolver
var queryMemberEvaluationsField = &graphql.Field{
    Type: graphql.NewList(queryMemberEvaluationObject),
    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        // ... // 
    },
}

// return object
var queryMemberEvaluationObject = graphql.NewObject(graphql.ObjectConfig{
    Name:   "memberEvaluation",
    Fields: graphql.Fields{
        // ... //
    },
}) 

```

Note the return object field name ("memberEvaluation") corresponds to the var name - 
query**MemberEvaluation**Object. 

This all seems very verbose, however there is so much potential for confusion and overlap 
between GraphQL fields names, returned object names, data types and data access functions 
that it is worth the pain. 

### Sample Queries

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


### Visual Representation

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




 
 