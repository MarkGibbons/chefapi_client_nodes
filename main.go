package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"github.com/gorilla/mux"
	"github.com/go-chef/chef"
	"github.com/MarkGibbons/chefapi_client"
	"github.com/MarkGibbons/chefapi_lib"
)

type restInfo struct {
        AuthUrl string
        Cert string
        Key string
        Port string
}

type NodeFilters struct {
	User string `json:"user,omitempty"`
	Organization string `json:"organization,omitempty"`
	NodeName string `json:"node_name,omitempty"`
}

type OrgNodes []NodeList

type NodeList struct {
	Organization string `json:"organization"`
	Nodes        []string `json:"nodes"`
}

var flags restInfo

func main() {
	flagInit()
	fmt.Printf("NODE FLAGS %+v\n", flags) //DEBUG

	r := mux.NewRouter()
	r.HandleFunc("/organizations", getNodes)
	r.HandleFunc("/organizations/{org}/nodes/{node}", singleNode)
	r.HandleFunc("/", defaultResp)
	// TODO: Verify that the request is authorized to call us
	// TODO: Use TLS
	log.Fatal(http.ListenAndServe("127.0.0.1:" + flags.Port, r))
	return
}

func getNodes( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	err := cleanInput(vars)
	if err != nil {
		inputerror(&w)
		return
	}
	// Get the filter information
	var filters NodeFilters
	err = json.NewDecoder(r.Body).Decode(&filters)
	// TODO: Allow for no filters
	// if err != nil {
	        // http.Error(w, err.Error(), http.StatusBadRequest)
		// return
	// }

	// Get a list of organizations to search for this request
	var orgList []string
	if filters.Organization == "" {
		orgList, err = allOrgs()
		fmt.Println("ORGLIST %+v ERR %+v\n", orgList, err)
		if err != nil {
			//TODO: deal with the error
		}
	} else {
		orgList[0] = filters.Organization
	}

	// Extract the node list
	orgNodes, err  :=  allNodes(orgList, filters)
	fmt.Println("ORGNODES %+v ERR %+v\n", orgNodes, err)
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

func singleNode( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	err := cleanInput(vars)
	if err != nil {
		inputerror(&w)
		return
	}

	switch r.Method {
	case "GET":
		// GET return single node information
		node, err := getNode(vars["org"],vars["node"])
		if err != nil {
			// TODO:
		}
		nodeJson,err := json.Marshal(node)
		if err != nil {
			// TODO:
		}
		w.WriteHeader(http.StatusOK)
		w.Write(nodeJson)
	case "POST":
		// PUT update a node
		var node chef.Node
		err = json.NewDecoder(r.Body).Decode(&node)
		if err != nil {
			// TODO:
			// handle the error
		}
		err = putNode(vars["org"],vars["node"],node)
		if err != nil {
			// TODO:
			// handle the error
		}
		w.WriteHeader(http.StatusOK)
	default:
	}
	return
}

func allNodes(orgs []string, filters NodeFilters) (orgnodes OrgNodes, err error) {
	for _,org := range orgs {
		client := chefapi_client.OrgClient(org)
		nodenames, err := listNodes(client, filters)
		if err != nil {
			// TODO: deal with an error
		}
		nodelist := NodeList{
			Organization: org,
			Nodes: nodenames,
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
	_,err = client.Nodes.Put(node)
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
	fmt.Println("in Allorgs")
        client := chefapi_client.Client()
	fmt.Println("passed client create")
        orgList, err  := client.Organizations.List()
	fmt.Println("passed orglist")
	orgNames = make([]string, 0, len(orgList))
	for k := range orgList {
		fmt.Printf("ADD TO ORGNAMES  K %+v", k)
		orgNames =  append(orgNames, k)
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
	return
}

func inputerror(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusBadRequest)
	(*w).Write([]byte(`{"message":"Bad url input value"}`))
}

func userAllowed(node string, user string) (authorized bool, err error) {
	authorized = false
	authurl := flags.AuthUrl  + "/auth/" + node + "/user/" + user
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
