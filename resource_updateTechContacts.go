package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUpdateTechContacts() *schema.Resource {
	return &schema.Resource{
		Create: resourceUpdateTechContactsCreate,
		Read:   resourceUpdateTechContactsRead,
		Update: resourceUpdateTechContactsUpdate,
		Delete: resourceUpdateTechContactsDelete,

		Schema: map[string]*schema.Schema{
			"technical_contact_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"phone_number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"help_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type technicalcontacts struct {
	TechnicalContactId string `json:"technicalContactId,omitempty"`
	PhoneNumber        string `json:"phoneNumber,omitempty"`
	HelpUrl            string `json:"helpUrl,omitempty"`
}

func resourceUpdateTechContactsCreate(d *schema.ResourceData, m interface{}) error {
	return resourceUpdateTechContactsUpdate(d, m)
}

func getContactsJson(body []byte) (*technicalcontacts, error) {
	var s = new(technicalcontacts)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error Unmarshalling Json Response:", err)
	}
	return s, err
}

func resourceUpdateTechContactsRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("https://%s-admin.%s/api/internal/enduser-support", m.(*Config).orgName, m.(*Config).domain)

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

	s, err := getContactsJson([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(fmt.Sprintf("%s.%s-TechnicalContacts", m.(*Config).orgName, m.(*Config).domain))
	d.Set("technical_contact_id", s.TechnicalContactId)
	d.Set("phone_number", s.PhoneNumber)
	d.Set("help_url", s.HelpUrl)
	return nil
}

func resourceUpdateTechContactsUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	url := fmt.Sprintf("https://%s-admin.%s/api/internal/enduser-support", m.(*Config).orgName, m.(*Config).domain)

	client := &http.Client{}

	org := technicalcontacts{
		TechnicalContactId: d.Get("technical_contact_id").(string),
		HelpUrl:            d.Get("help_url").(string),
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
		return fmt.Errorf("failed to Modify Interstitial Page Settings for url: %s, status: %s, req: %s", url, resp.Status, bytesJson)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getContactsJson([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.Partial(false)

	return resourceUpdateTechContactsRead(d, m)
}

func resourceUpdateTechContactsDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
