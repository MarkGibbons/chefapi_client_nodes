package main

import (
	 "fmt"
	 "net/http"
	 "net/http/httptest"
	 "reflect"
	 "testing"
	 "github.com/gorilla/mux"
	 "github.com/go-chef/chef"
)

const (
	userid = "tester"
	privateKey = "-----BEGIN RSA PRIVATE KEY-----"
)

var (
	chefmux *http.ServeMux
	server *httptest.Server
        client *chef.Client
)

// func getNodes( w http.ResponseWriter, r *http.Request) {
// func singleNode( w http.ResponseWriter, r *http.Request) {
// func allNodes(orgs []string, filters NodeFilters) (orgnodes OrgNodes, err error) {
// func getNode(organization string, nodename string) (node chef.Node, err error) {
// func putNode(organization string, nodenam string, node chef.Node) (err error) {
// func listNodes(client *chef.Client, filters NodeFilters) ([]string, error) {
// func allOrgs() (orgNames []string, err error) {
func TestAllOrgs(t *testing.T) {
        setup()
        defer teardown()

        chefmux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
                switch {
                case r.Method == "GET":
                        fmt.Fprintf(w, `{ "org_name1": "https://url/for/org_name1", "org_name2": "https://url/for/org_name2"}`)
                }
        })
	want := []string{"org_name1", "org_name2"}
	orgs, err := allOrgs()
        if err != nil {
                t.Errorf("AllOrgs unexpected error %v\n", err)
        }
	if !reflect.DeepEqual(orgs, want) {
                t.Errorf("Organizations.List returned %+v, want %+v", orgs, want)
        }
}

func TestCleanInput(t *testing.T) {
        expected_in := map[string]string{ "valid": "mynode" }
        err := cleanInput(expected_in)
        if err != nil {
                t.Errorf("Error cleaning: %+v Err: %+v\n", expected_in, err)
        }
        expected_in = map[string]string{ "invalid": "\nbounceit" }
        err = cleanInput(expected_in)
        if err == nil {
                t.Errorf("CleanInput did not receive expected error cleaning: %+v Err: %+v\n", expected_in, err)
        }
}

func  TestDefaultResp(t *testing.T) {
        req, err := http.NewRequest("GET", "/", nil)
        if err != nil {
                t.Fatal(err)
        }
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(defaultResp)
        handler.ServeHTTP(rr, req)

        // Check the status code and response body
        if status := rr.Code; status!= http.StatusBadRequest {
                t.Errorf("Status code is not expected. Got: %v want: %v\n", status, http.StatusBadRequest)
        }
        wantBody := `{"message":"GET /organizations/ORG/nodes\nGET /organizations/ORG/nodes/NODE\nPUT /organizations/ORG/nodes/NODE\nAre the only valid methods"}`
        if rr.Body.String() != wantBody {
                t.Errorf("AuthCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
        }
}

// TODO: func flagInit() (rest restInfo) {

// inputerror(w *http.ResponseWriter) {
func TestInputerror(t *testing.T) {
        // Check the status code and response body - invalid request invoked inputerror
        req, err := http.NewRequest("GET", "/organizations/other&org/nodes", nil)
        if err != nil {
                t.Fatal(err)
        }
        rr := httptest.NewRecorder()
        // Invoke the server
        newGetNodesServer().ServeHTTP(rr, req)
        // Check the status code and response body
        if status := rr.Code; status!= http.StatusBadRequest {
                t.Errorf("Get Nodes status code is not expected. Got: %v want: %v\n", status, http.StatusBadRequest)
        }
        wantBody := `{"message":"Bad url input value"}`
        if rr.Body.String() != wantBody {
                t.Errorf("Get Nodes unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
        }
}



func newGetNodesServer() http.Handler {
        r := mux.NewRouter()
        r.HandleFunc("/organizations/{org}/nodes", getNodes)
        return r
}

func newSingleNodeServer() http.Handler {
        r := mux.NewRouter()
        r.HandleFunc("/organizations/{org}/nodes/{node}", singleNode)
        return r
}

func setup() {
        chefmux = http.NewServeMux()
        server = httptest.NewServer(chefmux)
        client, _ = chef.NewClient(&chef.Config{
                Name:    userid,
                Key:     privateKey,
                BaseURL: server.URL,
        })
}

func teardown() {
         server.Close()
}
