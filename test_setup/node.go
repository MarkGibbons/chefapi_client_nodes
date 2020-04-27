package testapi

import (
	"fmt"
	"os"

	"github.com/go-chef/chef"
)

// Create nodes for organizations
// Create an assortment of nodes in each organization
func Node() {
	orgs := []string{"blue", "sales", "admin", "testorg", "mainapp"}
	node_prefixes := []string{"art", "ban", "red", "ska"}
	for i, org := range orgs {
		client := OrgClient(org)
		for j, pre := range node_prefixes {
			k := (i + 1) * (j + 1)
			nodename := fmt.Sprintf("%s%.4d", pre, k)
			node := chef.NewNode(nodename)
			node.RunList = []string{org + "::default"}
			_, err := client.Nodes.Post(node)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Couldn't create node node. ", err, nodename)
			}
		}
	}
}
