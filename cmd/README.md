## cmd/

Contains all of the executable packages:

* [pubmedr/](/cmd/pubmedr/README.md) - worker to fetch pubmed articles
* [mongr/](/cmd/mongr/README.md) - worker to sync data from MySQL to MongoDB
* [algr/](/cmd/algr/README.md) - worker to sync Algolia indexes
* [fixr/](/cmd/fixr/README.md) - utility to check and fix data
* [webd/](/cmd/webd/README.md) - web API, either REST or GraphQL<sup>1</sup>

<sup>1</sup>If env var `GRAPHQL_SERVER=true` then `webd` will start the GraphQL server. This is a workaround for the one web process limit on Heroku, allowing for the same repo to be pushed to two separate Heroku apps.

      




