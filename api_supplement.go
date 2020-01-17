package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/okta/okta-sdk-golang/okta"
)

// ApiSupplement not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type ApiSupplement struct {
	BaseURL         string
	Client          *http.Client
	Token           string
	RequestExecutor *okta.RequestExecutor
}

func (m *ApiSupplement) GetOrgById(id string) ([]byte, *http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/orgs/%s", m.BaseURL, id)
	return m.GetOrg(url)
}

func (m *ApiSupplement) GetOrg(url string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", m.Token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := m.RequestExecutor.DoWithRetries(okta.NewRequest(req), 0)
	if err != nil {
		return nil, res, err
	} else if res.StatusCode == http.StatusNotFound {
		return nil, res, nil
	} else if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("failed to get metadata for url: %s, status: %s", url, res.Status)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	return data, res, err
}
