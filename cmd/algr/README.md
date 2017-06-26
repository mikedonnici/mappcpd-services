# algr

Sync MappCPD data to Alogolia Indexes

[Algolia](https://www.algolia.com/) is used to create a better search experience in MappCPD for members, resources and modules.

At present the data flows as follows:

1. MappCPD admin app (ZF monolith) -> MySQL
2. [mongr](https://github.com/mappcpd/mongr) syncs data from MySQL -> MongoDB
3. algr syncs data from MongoDB -> Algolia indexes

Steps 2 & 3 both use the [api](https://github.com/mappcpd/api) to access data. The api itself 
can access both the MySQL and the MongoDB databases, depending on the endpoints. 

So mongr uses endpoints that get fetch data from MySQL, and save it to MongoDB, and algr leverages the faster api search 
endpoints that fetch data from MongoDB. 
  

## Env Vars
algr uses [envr](https://github.com/34South/envr) to verify required env vars, 
and will set them from a local **.env** file, if present. 
 
* ADMIN_USER='adminUserName'
* ADMIN_PASS='adminUserPass'
* ALG_APP_ID='AlgoliaAppID'
* ALG_API_KEY='AlgoliaAdminKey'
* BACK_DAYS=7
 
algr will sync documents with an _updatedAt_ date BACK_DAYS ago, or later.      




 