// A client library to access FORM3 API service

package accountsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const defaultEndpoint = "https://api.form3.tech"

type AccountAttributes struct {
	BankId       string `json:"bank_id"`
	BankIdCode   string `json:"bank_id_code"`
	BaseCurrency string `json:"base_currency"`
	Bic          string `json:"bic"`
	Country      string `json:"country"`
	//TODO add more https://api-docs.form3.tech/api.html#organisation-accounts-resource
}

type Account struct {
	Id             string            `json:"id"`
	OrganisationId string            `json:"organisation_id"`
	CreatedOn      string            `json:"created_on"`
	ModifiedOn     string            `json:"modified_on"`
	Version        int               `json:"version"`
	Type           string            `json:"type"`
	Attributes     AccountAttributes `json:"attributes"`
}

type AccountResponse struct {
	Account Account           `json:"data"`
	Links   map[string]string `json:"links"`
}

type AccountsResponse struct {
	Accounts []Account         `json:"data"`
	Links    map[string]string `json:"links"`
}

type createRequest struct {
	Account Account `json:"data"`
}

// An AccountsClient is an HTTP client.
// Default http.Client is used by default.
// Custom http/net client can be passed via composite literal.
type AccountsClient struct {
	httpClient http.Client
	endpoint   string
}

// Make HTTP request to a given route and
// return response content if http status is 2xx
func (c *AccountsClient) doRest(httpMethod string, route string, jsonData []byte) ([]byte, error) {
	if c.endpoint == "" {
		c.endpoint = defaultEndpoint
	}

	url := c.endpoint + route

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []byte{}, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, nil
	}

	// No worries about 3xx - net/http client does follow redirects by default
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return []byte{}, fmt.Errorf("%v != 2xx", resp.StatusCode)
	}

	//fmt.Println("response Body:", string(body))
	//fmt.Println("response Status:", resp.StatusCode)

	return body, nil
}

// Helper function for List() and ListPage()
func (c *AccountsClient) doList(url string) (AccountsResponse, error) {

	body, err := c.doRest("GET", url, nil)
	if err != nil {
		return AccountsResponse{}, err
	}

	var resp AccountsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return AccountsResponse{}, err
	}

	return resp, nil
}

// https://api-docs.form3.tech/api.html#organisation-accounts-create
func (c *AccountsClient) Create(id string, organisationId string, accountAttributes *AccountAttributes) (AccountResponse, error) {
	account := Account{
		Type:           "accounts",
		Id:             id,
		OrganisationId: organisationId,
		Attributes:     *accountAttributes,
	}

	req := createRequest{
		Account: account,
	}

	var jsonData []byte
	jsonData, err := json.Marshal(req)

	body, err := c.doRest("POST", "/v1/organisation/accounts", jsonData)
	if err != nil {
		return AccountResponse{}, err
	}

	var resp AccountResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return AccountResponse{}, err
	}

	return resp, nil
}

// https://api-docs.form3.tech/api.html#organisation-accounts-fetch
func (c *AccountsClient) Fetch(id string) (AccountResponse, error) {
	route := fmt.Sprintf("/v1/organisation/accounts/%s", id)

	body, err := c.doRest("GET", route, nil)
	if err != nil {
		return AccountResponse{}, err
	}

	var resp AccountResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return AccountResponse{}, err
	}

	return resp, nil
}

// https://api-docs.form3.tech/api.html#organisation-accounts-delete
func (c *AccountsClient) Delete(id string, version int) (bool, error) {

	route := fmt.Sprintf("/v1/organisation/accounts/%s?version=%d", id, version)
	_, err := c.doRest("DELETE", route, nil)

	if err != nil {
		return false, err
	}

	return true, nil
}

// https://api-docs.form3.tech/api.html#organisation-accounts-list
func (c *AccountsClient) List() (AccountsResponse, error) {
	return c.doList("/v1/organisation/accounts")
}

// https://api-docs.form3.tech/api.html#organisation-accounts-list
func (c *AccountsClient) ListPage(pageNumber int, pageSize int) (AccountsResponse, error) {
	return c.doList(fmt.Sprintf("/v1/organisation/accounts/?page[number]=%d&page[size]=%d", pageNumber, pageSize))
}
