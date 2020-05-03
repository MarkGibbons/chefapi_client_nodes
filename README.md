# chefapi_client_nodes

This code provides a go client interface to the chefapic to interact with nodes. Running this code provides a simpler API than the native chef REST API for use by web applications. See the chefapi_demo_server repository to see how this code was installed and started.

## Front End Endpoints used by web applications
-----------

### GET /orgnodes
===========================

### Request
Filter values may restrict the returned information to a specific owner of the
nodes and to a specific organization.  If not set nodes for all users and/or all
organizations will be returned.

#### Return
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
* 200 - List returned
* 400 - Invalid request was made
* 403 - Unauthorized request was made

### GET /orgnodes/ORG/nodes/NODE
===========================

#### Request

No request body is used

#### Return
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
* 200 - Node data returned
* 400 - Invalid request was made
* 403 - Unauthorized request was made

### PUT /orgnodes/ORG/nodes/NODE
===========================

#### Request
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

#### Return
 
No JSON body is returned.

Values
* 200 - Node data returned
* 400 - Invalid request was made
* 403 - Unauthorized request was made

## Back End Chef Infra Server Endpoints used
-----------

### GET /orgnodes/ORG/nodes
### GET /orgnodes/ORG/nodes/NODE
### PUT /orgnodes/ORG/nodes/NODE

## Links
-------
* https://blog.questionable.services/article/testing-http-handlers-go/
* https://github.com/quii/testing-gorillas
* https://godoc.org/github.com/gorilla/mux#SetURLVars
* https://github.com/gorilla/mux
* https://docs.chef.io/api_chef_server/
* https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
* https://blog.usejournal.com/authentication-in-golang-c0677bcce1a8 users and authentication scheme
* https://github.com/dgrijalva/jwt-go  json web token
* https://tutorialedge.net/golang/authenticating-golang-rest-api-with-jwts/ jwt basics
* https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/  jwt login, cookie, renewal
* https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS - cross origin resource sharing
* https://github.com/rs/cors - set up cors response in golang server
* https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body - clean up go http handlers
* https://stormpath.com/blog/where-to-store-your-jwts-cookies-vs-html5-web-storage - where to store JWTgit
