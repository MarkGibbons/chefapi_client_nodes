//
// Create test organizations
//
package testapi

import (
	"fmt"
	"os"

	"github.com/go-chef/chef"
)

func Organization() {
	// Create a client for access
	client := Client()

	// Organization tests
	orgs := []string{"blue", "sales", "admin", "testorg", "mainapp"}
	for _, org := range orgs {
		createOrganization(client, chef.Organization{Name: org, FullName: org})
	}

	orgList := listOrganizations(client)
	fmt.Println("List initial organizations", orgList)

}

// createOrganization uses the chef server api to create a single organization
func createOrganization(client *chef.Client, org chef.Organization) chef.OrganizationResult {
	orgResult, err := client.Organizations.Create(org)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating org:", org, err)
	}
	return orgResult
}

// deleteOrganization uses the chef server api to delete a single organization
func deleteOrganization(client *chef.Client, name string) error {
	err := client.Organizations.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", name, err)
	}
	return err
}

// getOrganization uses the chef server api to get information for a single organization
func getOrganization(client *chef.Client, name string) chef.Organization {
	orgList, err := client.Organizations.Get(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing org:", name, err)
	}
	return orgList
}

// listOrganizations uses the chef server api to list all organizations
func listOrganizations(client *chef.Client) map[string]string {
	orgList, err := client.Organizations.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing orgs:", err)
	}
	return orgList
}

// updateOrganization uses the chef server api to update information for a single organization
func updateOrganization(client *chef.Client, org chef.Organization) chef.Organization {
	org_update, err := client.Organizations.Update(org)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue updating org:", org, err)
	}
	return org_update
}
