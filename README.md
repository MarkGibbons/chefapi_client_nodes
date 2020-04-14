# chefapi_client_nodes

This code provides a go client interface to the chefapi to interact with nodes.
Running this code provides a simpler API than the native chef REST API for use by web applications.

# Front End Endpoints used by web applications
-----------

## GET /organizations/ORG/nodes
===========================

### Request
Filter values restrict the returned information to a specific owner of the 
nodes and to a specific organization.  If not set nodes for all users and/or all
organizations will be returned.

The request can specify filter values in a json body.
````json
{
  "user": "username",
  "organization" "orgname"
}
````


### Return
The body returned looks like this:
````json
{
  "organization": [
     "node1",
     "node2"
   ]
}
````

Values
200 - List returned
400 - Invalid request was made

## GET /organizations/ORG/nodes/NODE
===========================

### Request


### Return
The body returned looks like this:
````json
{
  "name": "node_name",
  "chef_environment": "_default",
  "run_list": [
    "recipe[recipe_name]"
  ]
  "json_class": "Chef::Node",
  "chef_type": "node",
  "automatic": { ... },
  "normal": { "tags": [ ] },
  "default": { },
  "override": { }
}
````

Values
200 - Node data returned
400 - Invalid request was made

## PUT /organizations/ORG/nodes/NODE
===========================

### Request
The request body looks like this:
````json
{
  "name": "node_name",
  "chef_environment": "_default",
  "run_list": [
    "recipe[recipe_name]"
  ]
  "json_class": "Chef::Node",
  "chef_type": "node",
  "automatic": { ... },
  "normal": { "tags": [ ] },
  "default": { },
  "override": { }
}
````

### Return
The body returned looks like this:
````json
{
}
````
Values
200 - Node data returned
400 - Invalid request was made

# Back End Chef Infra Server Endpoints used
-----------

## GET /organizations/ORG/nodes
## GET /organizations/ORG/nodes/NODE
## PUT /organizations/ORG/nodes/NODE

# Links
-------
https://blog.questionable.services/article/testing-http-handlers-go/
https://github.com/quii/testing-gorillas
https://godoc.org/github.com/gorilla/mux#SetURLVars
https://github.com/gorilla/mux
https://docs.chef.io/api_chef_server/
https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
https://blog.usejournal.com/authentication-in-golang-c0677bcce1a8 users and authentication scheme
https://github.com/dgrijalva/jwt-go  json web token
