package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			// "id": &schema.Schema{
			// 	Type: schema.TypeString,
			// 	Optional: true,
			// },
			"subdomain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"website": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"login": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobile_phone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"recovery_question": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"recovery_answer": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"secondary_email": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"secondary_image": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"error_redirect_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"interstitial_min_wait_time": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

type app struct {
	ERU  string `json:"errorRedirectUrl"`
	IMWT int    `json:"interstitialMinWaitTime"`
}

type attributes struct {
	SE bool `json:"secondaryEmail"`
	SI bool `json:"secondaryImage"`
}
type userAccount struct {
	A attributes `json:"attributes"`
}

type settings struct {
	A app         `json:"app"`
	U userAccount `json:"userAccount"`
}

type orgsettings struct {
	ID string   `json:"id"`
	SU string   `json:"subdomain"`
	N  string   `json:"name"`
	W  string   `json:"website"`
	S  settings `json:"settings"`
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	// subdomain := d.Get("subdomain").(string)
	// d.SetId(subdomain)
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	// id := d.Get("org_id").(string)

	url := fmt.Sprintf("https://%s.%s/api/v1/orgs/%s", m.(*Config).templateOrgName, m.(*Config).domain, m.(*Config).templateOrgName)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}
	// req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", d.Get("api_token")))
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

	s, err := getJsonResponse([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(s.ID)

	return nil
}

func getJsonResponse(body []byte) (*orgsettings, error) {
	var s = new(orgsettings)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error Unmarshalling Json Response:", err)
	}
	return s, err
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {

	d.Partial(true)

	// id := d.Get("org_id").(string)

	url := fmt.Sprintf("https://%s.%s/api/v1/orgs/%s", m.(*Config).templateOrgName, m.(*Config).domain, m.(*Config).templateOrgName)

	client := &http.Client{}

	appSettings := app{d.Get("error_redirect_url").(string), d.Get("interstitial_min_wait_time").(int)}

	attr := attributes{d.Get("secondary_email").(bool), d.Get("secondary_image").(bool)}

	ua := userAccount{attr}

	As := settings{appSettings, ua}

	org := orgsettings{"", d.Get("subdomain").(string), d.Get("name").(string), d.Get("website").(string), As}

	bytesJson, err := json.Marshal(org)

	if err != nil {
		fmt.Println("error:", err)
	}

	r := bytes.NewReader(bytesJson)

	req, err := http.NewRequest("PUT", url, r)

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
		return fmt.Errorf("failed to update Org Details for url: %s, status: %s, req: %s", url, resp.Status, bytesJson)
	}

	d.Partial(false)

	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
