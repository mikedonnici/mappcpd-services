## MappCPD Services Refactor

[![Build Status](https://travis-ci.org/mappcpd/web-services.svg?branch=master)](https://travis-ci.org/mappcpd/web-services)

Combine services projects into a single project structure based on Bill Kennedy's 
[package oriented design](https://www.goinggo.net/2017/02/package-oriented-design.html).


* [cmd/](/cmd/README.md) - all executable packages
  * [pubmedr/](/cmd/pubmedr/README.md) - pubmed article fetcher
  * [mongr/](/cmd/mongr/README.md) - syncs data from MySQL -> MongoDB
  * [pubmedr/](/cmd/algr/README.md) - sync Algolia search indexes
  * [webd/](/cmd/webd/README.md) - web server (api)
* [internal/](/internal/README.md) - in-house packages
* [vendor/](/vendor/README.md) - vendor packages

