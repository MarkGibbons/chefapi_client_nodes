package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/MarkGibbons/chefapi_client"
	"github.com/MarkGibbons/chefapi_lib"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chef/chef"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
)

type restInfo struct {
	AuthUrl string
	Cert    string
	Key     string
	Port    string
}

type NodeFilters struct {
	User         string `json:"user,omitempty"`
	Organization string `json:"organization,omitempty"`
	NodeName     string `json:"node,omitempty"`
}

type OrgNodes []NodeList

type NodeList struct {
	Organization string   `json:"organization"`
	Nodes        []string `json:"nodes"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var flags restInfo

var jwtKey = []byte("my_secret_key") // TODO: Parameter for production

func main() {
	flagInit()

	r := mux.NewRouter()
	r.HandleFunc("/orgnodes", getNodes)
	r.HandleFunc("/orgnodes/{org}/nodes/{node}", singleNode)
	r.HandleFunc("/", defaultResp)
	l, err := net.Listen("tcp4", ":"+flags.Port)
	if  err  != nil {
		panic(err.Error())
	}
	log.Fatal(http.ServeTLS(l, r, flags.Cert, flags.Key))
	return
}

func getNodes(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET NODES Method %+v\n", r.Method)
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	err := cleanInput(vars)
	if err != nil {
		inputerror(&w)
		return
	}
	// Get the filter information
	var filters NodeFilters
	// Get the filters from parameters
	userparm, ok := r.URL.Query()["user"]
	if ok {
		filters.User = userparm[0]
	}
	orgparm, ok := r.URL.Query()["organization"]
	if ok {
		filters.Organization = orgparm[0]
	}
	nodeparm, ok := r.URL.Query()["node"]
	if ok {
		filters.NodeName = nodeparm[0]
	}

	// Verify a logged in user made the request
	_, code := loggedIn(r)
	if code != -1 {
		fmt.Printf("Can't verify the user status: %+v\n", code)
		w.WriteHeader(code)
		return
	}

	// Get a list of organizations to search for this request
	fmt.Printf("Filters : %+v", filters)
	var orgList []string
	if filters.Organization == "" {
		orgList, err = allOrgs()
		if err != nil {
			//TODO: deal with the error
		}
	} else {
		orgList = append(orgList, filters.Organization)
	}

	// Extract the node list
	orgNodes, err := allNodes(orgList, filters)
	if err != nil {
		//TODO: Deal with the error
	}

	//  Handle the results and return the json body
	nodesJSON, err := json.Marshal(orgNodes)
	if err != nil {
		// TODO: deal with the error
	}
	w.WriteHeader(http.StatusOK)
	w.Write(nodesJSON)
	return
}

func singleNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	err := cleanInput(vars)
	if err != nil {
		inputerror(&w)
		return
	}

	switch r.Method {
	case "GET":
		// Verify a logged in user made the request
		_, code := loggedIn(r)
		if code != -1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// GET return single node information
		node, err := getNode(vars["org"], vars["node"])
		if err != nil {
			// TODO:
		}
		nodeJson, err := json.Marshal(node)
		if err != nil {
			// TODO:
		}
		w.WriteHeader(http.StatusOK)
		w.Write(nodeJson)
	case "PUT":
		// Verify a logged in user made the request
		user, code := loggedIn(r)
		if code != -1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// PUT update a node
		req, err := httputil.DumpRequest(r, true)
		fmt.Printf("PUT Request Body %+v\n %+v\n", string(req), err)
		var node chef.Node
		err = json.NewDecoder(r.Body).Decode(&node)
		if err != nil {
			fmt.Printf("POST JSON ERR: %+v\n", err)
			// TODO:
			// handle the error
		}

		// Verify the user is allowed to update this node
		userauth, err := userAllowed(node.Name, user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !userauth {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		fmt.Printf("UPDATE node %+v\n", node)
		err = putNode(vars["org"], vars["node"], node)
		if err != nil {
			fmt.Printf("Update JSON ERR: %+vi\n", err)
			// TODO:
			// handle the error
		}
		w.WriteHeader(http.StatusOK)
	default:
	}
	return
}

func allNodes(orgs []string, filters NodeFilters) (orgnodes OrgNodes, err error) {
	for _, org := range orgs {
		client := chefapi_client.OrgClient(org)
		nodenames, err := listNodes(client, filters)
		if err != nil {
			// TODO: deal with an error
		}
		nodelist := NodeList{
			Organization: org,
			Nodes:        nodenames,
		}
		orgnodes = append(orgnodes, nodelist)
	}
	return
}

func getNode(organization string, nodename string) (node chef.Node, err error) {
	client := chefapi_client.OrgClient(organization)
	node, err = client.Nodes.Get(nodename)
	fmt.Printf("GETNODE %+v ERR %+v NODE %+v ORG %+v\n", node, err, nodename, organization)
	return
}

func putNode(organization string, nodenam string, node chef.Node) (err error) {
	client := chefapi_client.OrgClient(organization)
	nodereturn, err := client.Nodes.Put(node)
	fmt.Printf("Return from put node %+v\n", nodereturn)
	return
}

func listNodes(client *chef.Client, filters NodeFilters) (nodeNames []string, err error) {
	nodeList, err := client.Nodes.List()
	fmt.Printf("NODELIST %+v CLIENT %+v  ERR %+v\n", nodeList, client, err)
	if err != nil {
		// TODO: Deal with errors
		return
	}
	nodeNames = make([]string, 0, len(nodeList))
	nameMatcher, err := regexp.Compile(filters.NodeName)
	if err != nil {
		err = errors.New("Invalid regular expression for the node name filter")
		return
	}
	for node := range nodeList {
		// apply the name filter
		if !nameMatcher.Match([]byte(node)) {
			continue
		}
		// check owner specified and filter by ownership
		if filters.User != "" {
			allowed, err := userAllowed(node, filters.User)
			if err != nil || !allowed {
				continue
			}
		}
		nodeNames = append(nodeNames, node)
	}
	return nodeNames, err
}

func allOrgs() (orgNames []string, err error) {
	client := chefapi_client.Client()
	orgList, err := client.Organizations.List()
	orgNames = make([]string, 0, len(orgList))
	for k := range orgList {
		orgNames = append(orgNames, k)
	}
	return
}

func cleanInput(vars map[string]string) (err error) {
	for _, value := range vars {
		matched, _ := regexp.MatchString("^[[:word:]]+$", value)
		if !matched {
			err = errors.New("Invalid value in the URI")
			break
		}
	}
	return
}

func defaultResp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"message":"GET /organizations/ORG/nodes\nGET /organizations/ORG/nodes/NODE\nPUT /organizations/ORG/nodes/NODE\nAre the only valid methods"}`))
}

func flagInit() {
	restcert := flag.String("restcert", "", "Rest Certificate File")
	restkey := flag.String("restkey", "", "Rest Key File")
	restport := flag.String("restport", "8111", "Rest interface https port")
	authurl := flag.String("authurl", "", "Node authorization service url")
	flag.Parse()
	flags.AuthUrl = *authurl
	flags.Cert = *restcert
	flags.Key = *restkey
	flags.Port = *restport
	fmt.Printf("Flags used %+v\n", flags)
	return
}

func inputerror(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusBadRequest)
	(*w).Write([]byte(`{"message":"Bad url input value"}`))
}

func userAllowed(node string, user string) (authorized bool, err error) {
	authorized = false
	authurl := flags.AuthUrl + "/auth/" + node + "/user/" + user
	resp, err := http.Get(authurl)
	if err != nil {
		return
	}
	var auth chefapi_lib.Auth
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		return
	}
	authorized = auth.Auth
	return
}

// loggedIn verifies the JWT and extracts the user name
func loggedIn(r *http.Request) (user string, code int) {
	code = -1
	reqToken := r.Header.Get("Authorization")
	fmt.Printf("REQTOKEN %+v\n", reqToken)
	splitToken := strings.Split(reqToken, "Bearer")
	// Verify index before using
	if len(splitToken) != 2 {
		code = http.StatusBadRequest
		return
	}
	tknStr := strings.TrimSpace(splitToken[1])
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			code = http.StatusUnauthorized
			return
		}
		code = http.StatusBadRequest
		return
	}
	if !tkn.Valid {
		code = http.StatusUnauthorized
		return
	}
	user = claims.Username
	return
}
