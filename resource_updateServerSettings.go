package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceupdateServerSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceupdateServerSettingsCreate,
		Read:   resourceupdateServerSettingsRead,
		Update: resourceupdateServerSettingsUpdate,
		Delete: resourceupdateServerSettingsDelete,

		Schema: map[string]*schema.Schema{
			"subdomain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
			"app_error_redirect_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"portal_error_redirect_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"portal_signout_url": &schema.Schema{
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

func resourceupdateServerSettingsCreate(d *schema.ResourceData, m interface{}) error {
	// subdomain := d.Get("subdomain").(string)
	// d.SetId(subdomain)
	return resourceupdateServerSettingsUpdate(d, m)
}

func resourceupdateServerSettingsRead(d *schema.ResourceData, m interface{}) error {
	// id := d.Get("org_id").(string)

	url := fmt.Sprintf("https://%s.%s/api/v1/orgs/%s", m.(*Config).orgName, m.(*Config).domain, m.(*Config).orgName)

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

	s, err := getJsonResponse([]byte(body))
	if err != nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(s.ID)
	d.Set("website", s.Website)
	d.Set("name", s.Name)
	d.Set("subdomain", s.Subdomain)
	d.Set("secondary_email", s.Settings.UserAccount.Attributes.SecondaryEmail)
	d.Set("secondary_image", s.Settings.UserAccount.Attributes.SecondaryImage)
	d.Set("app_error_redirect_url", s.Settings.App.ErrorRedirectURL)
	d.Set("portal_error_redirect_url", s.Settings.Portal.ErrorRedirectURL)
	d.Set("portal_signout_url", s.Settings.Portal.SignOutURL)
	d.Set("interstitial_min_wait_time", s.Settings.App.InterstitialMinWaitTime)

	return nil
}

func getJsonResponse(body []byte) (orgsettings, error) {
	var s orgsettings
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error Unmarshalling Json Response:", err)
	}
	return s, err
}

func resourceupdateServerSettingsUpdate(d *schema.ResourceData, m interface{}) error {

	d.Partial(true)

	url := fmt.Sprintf("https://%s.%s/api/v1/orgs/%s", m.(*Config).orgName, m.(*Config).domain, m.(*Config).orgName)

	client := &http.Client{}

	org := orgsettings{
		Name:    d.Get("name").(string),
		Website: d.Get("website").(string),
		Settings: &Settings{
			App: &App{
				ErrorRedirectURL:        d.Get("app_error_redirect_url").(string),
				InterstitialMinWaitTime: d.Get("interstitial_min_wait_time").(int),
			},
			UserAccount: &UserAccount{
				Attributes: &Attributes{
					SecondaryEmail: d.Get("secondary_email").(bool),
					SecondaryImage: d.Get("secondary_image").(bool),
				},
			},
			Portal: &Portal{
				ErrorRedirectURL: d.Get("portal_error_redirect_url").(string),
				SignOutURL:       d.Get("portal_signout_url").(string),
			},
		},
	}

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getJsonResponse([]byte(body))
	if err != nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(s.ID)
	d.Set("website", s.Website)
	d.Set("name", s.Name)
	d.Set("subdomain", s.Subdomain)
	d.Set("secondary_email", s.Settings.UserAccount.Attributes.SecondaryEmail)
	d.Set("secondary_image", s.Settings.UserAccount.Attributes.SecondaryImage)
	d.Set("app_error_redirect_url", s.Settings.App.ErrorRedirectURL)
	d.Set("portal_error_redirect_url", s.Settings.Portal.ErrorRedirectURL)
	d.Set("portal_signout_url", s.Settings.Portal.SignOutURL)
	d.Set("interstitial_min_wait_time", s.Settings.App.InterstitialMinWaitTime)

	d.Partial(false)

	return nil

}

func resourceupdateServerSettingsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
