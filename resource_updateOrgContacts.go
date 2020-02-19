package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceupdateOrgContacts() *schema.Resource {
	return &schema.Resource{
		Create: resourceupdateOrgContactsCreate,
		Read:   resourceupdateOrgContactsRead,
		Update: resourceupdateOrgContactsUpdate,
		Delete: resourceupdateOrgContactsDelete,

		Schema: map[string]*schema.Schema{
			"street_address_1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"street_address_2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"city": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zip_code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"country": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mobile_phone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"technical_phone_number": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceupdateOrgContactsCreate(d *schema.ResourceData, m interface{}) error {
	// subdomain := d.Get("subdomain").(string)
	// d.SetId(subdomain)
	return resourceupdateOrgContactsUpdate(d, m)
}

func resourceupdateOrgContactsRead(d *schema.ResourceData, m interface{}) error {
	// id := d.Get("org_id").(string)

	url := fmt.Sprintf("https://%s.%s/api/v1/orgs/%s/contacts", m.(*Config).orgName, m.(*Config).domain, m.(*Config).orgName)

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

	s, err := getOrgContactsResponse([]byte(body))
	if err != nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(fmt.Sprintf("%s.%s-orgcontactsettings", m.(*Config).orgName, m.(*Config).domain))
	d.Set("street_address_1", s.Orgcontact.Steet_Address[0])
	d.Set("street_address_2", s.Orgcontact.Steet_Address[1])
	d.Set("city", s.Orgcontact.City)
	d.Set("state", s.Orgcontact.State)
	d.Set("zip_code", s.Orgcontact.ZipCode)
	d.Set("country", s.Orgcontact.Country)
	d.Set("mobile_phone", s.Orgcontact.Phone)
	d.Set("technical_phone_number", s.Support.Tech_Contact_Phone)

	return nil
}

type Support struct {
	Tech_Contact_Phone string `json:"phoneNumber,omitempty"`
}

type Orgcontact struct {
	City          string   `json:"locality,omitempty"`
	State         string   `json:"region,omitempty"`
	ZipCode       string   `json:"postalCode,omitempty"`
	Country       string   `json:"country,omitempty"`
	Phone         string   `json:"phoneNumber,omitempty"`
	Steet_Address []string `json:"streetAddress,omitempty"`
}

type orgcontactsjson struct {
	Orgcontact *Orgcontact `json:"organization"`
	Support    *Support    `json:"support"`
}

func getOrgContactsResponse(body []byte) (orgcontactsjson, error) {
	var s orgcontactsjson
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error Unmarshalling Json Response:", err)
	}
	return s, err
}

func resourceupdateOrgContactsUpdate(d *schema.ResourceData, m interface{}) error {

	d.Partial(true)

	url := fmt.Sprintf("https://%s.%s/api/v1/orgs/%s/contacts", m.(*Config).orgName, m.(*Config).domain, m.(*Config).orgName)

	client := &http.Client{}

	org := orgcontactsjson{
		Orgcontact: &Orgcontact{
			City:    d.Get("city").(string),
			State:   d.Get("state").(string),
			ZipCode: d.Get("zip_code").(string),
			Country: d.Get("country").(string),
			Phone:   d.Get("mobile_phone").(string),
			Steet_Address: []string{
				d.Get("street_address_1").(string),
				d.Get("street_address_2").(string),
			},
		},
		Support: &Support{
			Tech_Contact_Phone: d.Get("technical_phone_number").(string),
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
		return fmt.Errorf("failed to update Org  Contact Details for url: %s, status: %s, req: %s", url, resp.Status, bytesJson)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getOrgContactsResponse([]byte(body))
	if err != nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(fmt.Sprintf("%s.%s-orgcontactsettings", m.(*Config).orgName, m.(*Config).domain))
	d.Set("street_address_1", s.Orgcontact.Steet_Address[0])
	d.Set("street_address_2", s.Orgcontact.Steet_Address[1])
	d.Set("city", s.Orgcontact.City)
	d.Set("state", s.Orgcontact.State)
	d.Set("zip_code", s.Orgcontact.ZipCode)
	d.Set("country", s.Orgcontact.Country)
	d.Set("mobile_phone", s.Orgcontact.Phone)
	d.Set("technical_phone_number", s.Support.Tech_Contact_Phone)

	d.Partial(false)

	return nil

}

func resourceupdateOrgContactsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
