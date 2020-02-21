package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCreateAPITokens() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateAPITokensCreate,
		Read:   resourceCreateAPITokensRead,
		Update: resourceCreateAPITokensUpdate,
		Delete: resourceCreateAPITokensDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"apikey": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type apiTokens struct {
	TokenType  string `json:"tokenType,omitempty"`
	ClientName string `json:"clientName,omitempty"`
	Status     string `json:"status,omitempty"`
	Name       string `json:"name,omitempty"`
	Id         string `json:"id,omitempty"`
	Token      string `json:"token,omitempty"`
}

func resourceCreateAPITokensCreate(d *schema.ResourceData, m interface{}) error {
	return resourceCreateAPITokensUpdate(d, m)
}

func getapiTokensJson(body []byte) (*apiTokens, error) {
	var s = new(apiTokens)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error Unmarshalling Json Response:", err)
	}
	return s, err
}

func resourceCreateAPITokensRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("https://%s-admin.%s/api/internal/tokens/%s", m.(*Config).orgName, m.(*Config).domain, d.Id())

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", m.(*Config).apiToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		return nil
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get Org Details for url: %s, status: %s", url, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getapiTokensJson([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(s.Id)
	d.Set("name", s.Name)
	d.Set("token_type", s.TokenType)
	d.Set("status", s.Status)
	d.Set("client_name", s.ClientName)
	return nil
}

func resourceCreateAPITokensUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	url := fmt.Sprintf("https://%s-admin.%s/api/internal/tokens", m.(*Config).orgName, m.(*Config).domain)

	client := &http.Client{}

	org := apiTokens{
		Name: d.Get("name").(string),
	}

	bytesJson, err := json.Marshal(org)

	if err != nil {
		fmt.Println("error:", err)
	}

	r := bytes.NewReader(bytesJson)

	req, err := http.NewRequest("POST", url, r)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", m.(*Config).apiToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		return nil
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to Create API Tokens : %s, status: %s", url, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getapiTokensJson([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(s.Id)
	d.Set("name", s.Name)
	d.Set("token_type", s.TokenType)
	d.Set("status", s.Status)
	d.Set("client_name", s.ClientName)
	d.Set("apikey", s.Token)

	d.Partial(false)

	return nil
}

func resourceCreateAPITokensDelete(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("https://%s-admin.%s/api/internal/tokens/%s/revoke", m.(*Config).orgName, m.(*Config).domain, d.Id())

	client := &http.Client{}

	org := apiTokens{
		Name: d.Get("name").(string),
	}

	bytesJson, err := json.Marshal(org)

	if err != nil {
		fmt.Println("error:", err)
	}

	r := bytes.NewReader(bytesJson)

	req, err := http.NewRequest("POST", url, r)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", m.(*Config).apiToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		return nil
	} else if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to destroy API Token: %s, status: %s", url, resp.Status)
	}

	d.SetId("")
	return nil
}
