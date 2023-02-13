// A client library to access FORM3 API service

package accountsclient

import (
	"fmt"
	"testing"
)

// Sample ids
const accId = "0d27e265-9605-4b4b-a0e5-3003ea9cc4db"
const orgId = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73b"

// Testing endpoint
const endpoint = "http://accountapi:8080"

// Sample account attributes
var accAttributes = AccountAttributes{
	Country:      "GB",
	BaseCurrency: "GBP",
	BankId:       "400300",
	BankIdCode:   "GBDSC",
	Bic:          "NWBKGB22",
}

// Helper that cleans the testing account
// to let test cases to run in defined state
func setUp() AccountsClient {
	client := AccountsClient{endpoint: endpoint}
	client.Delete(accId, 0)
	return client
}

// List accounts for sake of asserting there are @expected accounts
func assertXAccounts(t *testing.T, client *AccountsClient, expected int) {
	resp, err := client.List()

	if err != nil {
		t.Errorf("Error listing accounts: %v", err)
		return
	}

	if len(resp.Accounts) != expected {
		t.Errorf("List returned %d != %d", len(resp.Accounts), expected)
		return
	}
}

// Helper - generate stable nth uuid based on accId
func nthAccId(n int) string {
	return fmt.Sprintf("%d", n) + accId[1:]
}

// Basic test of all API operations (Create, Fetch, List, Delete)
// Detailed test for each operation in separate functions.
func TestAllOperations(t *testing.T) {
	client := setUp()

	// Make sure there are no accounts
	assertXAccounts(t, &client, 0)

	// 1/5 Create account
	{
		_, err := client.Create(accId, orgId, &accAttributes)
		if err != nil {
			t.Errorf("Error creating account: %v", err)
			return
		}
	}

	// 2/5 Fetch the account
	{
		_, err := client.Fetch(accId)
		if err != nil {
			t.Errorf("Error fetching account: %v", err)
			return
		}
	}

	// 3/5 List accounts to check if the account is returned
	{
		resp, err := client.List()

		if err != nil {
			t.Errorf("Error listing accounts: %v", err)
			return
		}

		if len(resp.Accounts) != 1 {
			t.Errorf("List returned %d != 1", len(resp.Accounts))
			return
		}
	}

	// 4/5 List using pagination
	{
		resp, err := client.ListPage(0, 10)
		if err != nil {
			t.Errorf("Error listing accounts: %v", err)
			return
		}

		if len(resp.Accounts) != 1 {
			t.Errorf("List returned %d != 1", len(resp.Accounts))
			return
		}
	}

	// 5/5 Delete account
	{
		resp, err := client.Delete(accId, 0)
		if err != nil {
			t.Errorf("Error listing accounts: %v", err)
			return
		}

		if resp != true {
			t.Errorf("List returned %v != true", resp)
			return
		}
	}

	// Make sure there are no accounts
	assertXAccounts(t, &client, 0)
}

func TestInvalidEndpoint(t *testing.T) {

	client := AccountsClient{endpoint: "http://localhost:1234"}

	_, err := client.Create(accId, orgId, &accAttributes)

	if err == nil {
		t.Errorf("Error connecting to fake endpoint was expected")
		return
	}
}

func TestCreate(t *testing.T) {
	client := setUp()

	resp, _ := client.Create(accId, orgId, &accAttributes)
	the_account := resp.Account

	if the_account.Id != accId {
		t.Errorf("Account ID not saved")
	}

	if the_account.OrganisationId != orgId {
		t.Errorf("Account OrganisationId ID not saved")
	}

	//TODO: add more validation, however it is
	// supposed to be covered in API itself tests
}

func TestFetch(t *testing.T) {
	client := setUp()

	resp, _ := client.Create(accId, orgId, &accAttributes)
	the_account := resp.Account

	if the_account.Id != accId {
		t.Errorf("Account ID not saved")
	}

	if the_account.OrganisationId != orgId {
		t.Errorf("Account OrganisationId ID not saved")
	}

	//TODO: add more validation, however it is
	// supposed to be covered in API tests
}

func TestDelete(t *testing.T) {
	client := setUp()

	resp, _ := client.Create(accId, orgId, &accAttributes)
	the_account := resp.Account

	assertXAccounts(t, &client, 1)
	client.Delete(the_account.Id, the_account.Version)
	assertXAccounts(t, &client, 0)
}

func TestList(t *testing.T) {
	client := setUp()

	// Insert 10 accounts
	for i := 0; i < 10; i++ {
		client.Create(nthAccId(i), orgId, &accAttributes)
	}

	resp, _ := client.List()

	// Clean
	for i := 0; i < 10; i++ {
		client.Delete(nthAccId(i), resp.Accounts[i].Version)
	}

	if len(resp.Accounts) != 10 {
		t.Errorf("List returned %d != 10", len(resp.Accounts))
	}

	// Assert content
	for i := 0; i < 10; i++ {
		id := nthAccId(i)
		if resp.Accounts[i].Id != id {
			t.Errorf("List returned ID %v != %v", resp.Accounts[i].Id, id)
		}

		if resp.Accounts[i].OrganisationId != orgId {
			t.Errorf("List returned Organisation ID %v != %v", resp.Accounts[i].OrganisationId, orgId)
		}
	}
}

func TestListPage(t *testing.T) {
	client := setUp()

	// Insert 10 accounts
	for i := 0; i < 10; i++ {
		client.Create(nthAccId(i), orgId, &accAttributes)
	}

	// Request 5 accounts from page number 1 (2nd page)
	resp, _ := client.ListPage(1, 5)

	// Clean
	for i := 0; i < 10; i++ {
		client.Delete(nthAccId(i), 0)
	}

	// Assert content
	for i := 5; i < 10; i++ {
		id := nthAccId(i)
		if resp.Accounts[i-5].Id != id {
			t.Errorf("List returned %v != %v", resp.Accounts[i-5].Id, id)
		}
	}
}
