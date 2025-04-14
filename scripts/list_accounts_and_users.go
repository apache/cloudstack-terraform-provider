package main

import (
	"fmt"
	"log"
	"os"

	"github.com/apache/cloudstack-go/v2/cloudstack"
)

func main() {
	// Create a new CloudStack client
	apiURL := os.Getenv("CLOUDSTACK_API_URL")
	apiKey := os.Getenv("CLOUDSTACK_API_KEY")
	secretKey := os.Getenv("CLOUDSTACK_SECRET_KEY")

	if apiURL == "" || apiKey == "" || secretKey == "" {
		log.Fatal("CLOUDSTACK_API_URL, CLOUDSTACK_API_KEY, and CLOUDSTACK_SECRET_KEY must be set")
	}

	cs := cloudstack.NewClient(apiURL, apiKey, secretKey, false)

	// List accounts
	fmt.Println("=== Accounts ===")
	p := cs.Account.NewListAccountsParams()
	accounts, err := cs.Account.ListAccounts(p)
	if err != nil {
		log.Fatalf("Error listing accounts: %v", err)
	}

	for _, account := range accounts.Accounts {
		fmt.Printf("Account: %s, ID: %s\n", account.Name, account.Id)
	}

	// List users
	fmt.Println("\n=== Users ===")
	u := cs.User.NewListUsersParams()
	users, err := cs.User.ListUsers(u)
	if err != nil {
		log.Fatalf("Error listing users: %v", err)
	}

	for _, user := range users.Users {
		fmt.Printf("User: %s, ID: %s, Account: %s\n", user.Username, user.Id, user.Account)
	}
}
