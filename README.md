## MappCPD Services Refactor

Combine services projects into a single project structure based on Bill Kennedy's 
[package oriented design](https://www.goinggo.net/2017/02/package-oriented-design.html).


* [cmd/](/cmd/README.md) - all executable packages
  * [webd/](/cmd/webd/README.md) - web server (api)
* [internal/](/internal/README.md) - in-house packages
