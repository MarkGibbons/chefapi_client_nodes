package testapi

import (
	"fmt"
	"os"

	"github.com/go-chef/chef"
)

// Create users
// Add the users to the organizations
func User() {
	orgs := []string{"blue", "sales", "admin", "testorg", "mainapp"}
	users := []string{"arline", "doan", "fred", "sue"}
	client := Client()
	for _, user := range users {
		usr := chef.User{
			UserName:    user,
			Email:       user + "@domain.io",
			FirstName:   user,
			LastName:    "fullname",
			DisplayName: user + "Fullname",
			Password:    "Logn12ComplexPwd#",
			CreateKey:   true,
		}
		_, err := client.Users.Create(usr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't create user %+v Err: %+v\n", user, err)
		}
	}

	for _, org := range orgs {
		client := OrgClient(org)
		for _, user := range users {
			adduser := chef.AddNow{Username: user}
			err := client.Associations.Add(adduser)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't add user %+v to org %+v Err: %+v\n", user, org, err)
			}
			group, err := client.Groups.Get("admins")
			fmt.Printf("Group list  %+v in org %+v\n", group, org)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't get the admins group for org %+v Err: %+v\n", org, err)
			}
			group.Actors = append(group.Actors, "sue")
			group.Users = append(group.Users, "sue")
			groupupd, err := client.Groups.Update(group)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't add sue to the admins group for org %+v Err: %+v\n", org, err)
			}
			fmt.Printf("Group result %+v in org %+v\n", groupupd, org)
		}
	}
	client = OrgClient("blue")
	group := chef.Group{}
	group.Name = "blueusers"
	group.GroupName = "blueusers"
	group.Actors = []string{"sue"}
	group.Users = []string{"sue"}
	groupres, err := client.Groups.Create(group)
	fmt.Printf("Blueusers create %+v %+v\n", groupres, err)
	group.Users = []string{"sue", "doan"}
	groupupd, err := client.Groups.Update(group)
	fmt.Printf("Blueusers %+v %+v\n", groupupd, err)
	groupout, err := client.Groups.Get("blueusers")
	fmt.Printf("Blueusers %+v %+v\n", groupout, err)
}
