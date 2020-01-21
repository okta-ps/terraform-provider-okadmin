package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceInterstitialSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceInterstitialSettingsCreate,
		Read:   resourceInterstitialSettingsRead,
		Update: resourceInterstitialSettingsUpdate,
		Delete: resourceInterstitialSettingsDelete,

		Schema: map[string]*schema.Schema{
			"interstitial_page_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

type interstitialsettings struct {
	OIE bool `json:"oktaInterstitialEnabled"`
}

func resourceInterstitialSettingsCreate(d *schema.ResourceData, m interface{}) error {
	// oktaInterstitialPageId := d.Get("interstitial_page_id").(string)
	// d.SetId(oktaInterstitialPageId)
	// return resourceInterstitialSettingsRead(d, m)
	return resourceInterstitialSettingsRead(d, m)
}

func getJsonResponseInt(body []byte) (*interstitialsettings, error) {
	var s = new(interstitialsettings)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error Unmarshalling Json Response:", err)
	}
	return s, err
}

func resourceInterstitialSettingsRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("https://%s-admin.%s/api/internal/v1/oktaInterstitial/settings", m.(*Config).templateOrgName, m.(*Config).domain)

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

	s, err := getJsonResponseInt([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	// d.SetId(d.Id()
	d.Set("interstitial_page_enabled", s.OIE)

	return nil
}

func resourceInterstitialSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	url := fmt.Sprintf("https://%s-admin.%s/api/internal/v1/oktaInterstitial/settings", m.(*Config).templateOrgName, m.(*Config).domain)

	client := &http.Client{}

	org := interstitialsettings{d.Get("interstitial_page_enabled").(bool)}

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
		return fmt.Errorf("failed to Modify Interstitial Page Settings for url: %s, status: %s, req: %s", url, resp.Status, bytesJson)
	}

	d.Partial(false)

	return resourceInterstitialSettingsRead(d, m)
}

func resourceInterstitialSettingsDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
