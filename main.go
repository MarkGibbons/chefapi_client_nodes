package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/MarkGibbons/chefapi_client"
	"github.com/MarkGibbons/chefapi_lib"
	"github.com/go-chef/chef"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"regexp"
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

var flags restInfo

func main() {
	flagInit()

	r := mux.NewRouter()
	r.HandleFunc("/orgnodes", getNodes)
	r.HandleFunc("/orgnodes/{org}/nodes/{node}", singleNode)
	l, err := net.Listen("tcp4", ":"+flags.Port)
	if err != nil {
		panic(err.Error())
	}
	log.Fatal(http.ServeTLS(l, r, flags.Cert, flags.Key))
	return
}

func getNodes(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	err := chefapi_lib.CleanInput(vars)
	if err != nil {
		chefapi_lib.InputError(&w)
		return
	}

	// Get the filters from parameters
	var filters NodeFilters
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
	_, code := chefapi_lib.LoggedIn(r)
	if code != -1 {
		fmt.Printf("Can't verify the user status: %+v\n", code)
		w.WriteHeader(code)
		return
	}

	// Get a list of organizations to search for this request
	var orgList []string
	if filters.Organization == "" {
		orgList, err = chefapi_lib.AllOrgs()
		if err != nil {
			msg, code := chefapi_lib.ChefStatus(err)
			http.Error(w, msg, code)
			return
		}
	} else {
		orgList = append(orgList, filters.Organization)
	}

	// Extract the node list
	orgNodes, err := allNodes(orgList, filters)
	if err != nil {
		msg, code := chefapi_lib.ChefStatus(err)
		http.Error(w, msg, code)
		return
	}

	//  Handle the results and return the json body
	nodesJSON, err := json.Marshal(orgNodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(nodesJSON)
	return
}

func singleNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	err := chefapi_lib.CleanInput(vars)
	if err != nil {
		chefapi_lib.InputError(&w)
		return
	}

	switch r.Method {
	case "GET":
		// Verify a logged in user made the request
		_, code := chefapi_lib.LoggedIn(r)
		if code != -1 {
			http.Error(w,"User is not logged in", http.StatusUnauthorized)
			return
		}

		// GET return single node information
		node, err := getNode(vars["org"], vars["node"])
		if err != nil {
			msg, code := chefapi_lib.ChefStatus(err)
			http.Error(w, msg, code)
			return
		}
		nodeJson, err := json.Marshal(node)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(nodeJson)
	case "PUT":
		// Verify a logged in user made the request
		user, code := chefapi_lib.LoggedIn(r)
		if code != -1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// PUT update a node
		var node chef.Node
		err = json.NewDecoder(r.Body).Decode(&node)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Verify the user is allowed to update this node
		userauth, err := userAllowed(node.Name, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !userauth {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		fmt.Printf("UPDATE node %+v\n", node)
		err = putNode(vars["org"], vars["node"], node)
		if err != nil {
			msg, code := chefapi_lib.ChefStatus(err)
			http.Error(w, msg, code)
			return
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
			continue
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
	return
}

func putNode(organization string, nodenam string, node chef.Node) (err error) {
	client := chefapi_client.OrgClient(organization)
	_, err = client.Nodes.Put(node)
	return
}

func listNodes(client *chef.Client, filters NodeFilters) (nodeNames []string, err error) {
	nodeList, err := client.Nodes.List()
	if err != nil {
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

func userAllowed(node string, user string) (authorized bool, err error) {
	authorized = false
	authurl := flags.AuthUrl + "/auth/" + user + "/node/" + node
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
