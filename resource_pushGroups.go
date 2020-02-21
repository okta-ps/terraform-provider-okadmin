package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcepushGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourcepushGroupsCreate,
		Read:   resourcepushGroupsRead,
		Update: resourcepushGroupsUpdate,
		Delete: resourcepushGroupsDelete,

		Schema: map[string]*schema.Schema{
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"org_to_org_instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type pushgroupsreq struct {
	Status      string `json:"status,omitempty"`
	UserGroupId string `json:"userGroupId,omitempty"`
}

type pushgroupsres struct {
	Status    string `json:"status,omitempty"`
	MappingId string `json:"mappingId,omitempty"`
}

func resourcepushGroupsCreate(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("https://%s-admin.%s/api/internal/instance/%s/grouppush", m.(*Config).orgName, m.(*Config).domain, d.Get("org_to_org_instance_id"))

	org := pushgroupsreq{
		Status:      d.Get("status").(string),
		UserGroupId: d.Get("user_group_id").(string),
	}

	client := &http.Client{}

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
		return fmt.Errorf("failed to Push groups for Org2Org app for url: %s, status: %s, req: %s", url, resp.Status, bytesJson)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getPushGroupsResJson([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(s.MappingId)
	d.Set("status", s.Status)

	return nil

}

func getPushGroupsReqJson(body []byte) (*pushgroupsreq, error) {
	var s = new(pushgroupsreq)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error Unmarshalling Json Response:", err)
	}
	return s, err
}

func getPushGroupsResJson(body []byte) (*pushgroupsres, error) {
	var s = new(pushgroupsres)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error Unmarshalling Json Response:", err)
	}
	return s, err
}

func resourcepushGroupsRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("https://%s-admin.%s/api/internal/instance/%s/grouppush/%s", m.(*Config).orgName, m.(*Config).domain, d.Get("org_to_org_instance_id").(string), d.Id())

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

	s, err := getPushGroupsResJson([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(s.MappingId)
	d.Set("status", s.Status)
	return nil
}

func resourcepushGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	url := fmt.Sprintf("https://%s-admin.%s/api/internal/instance/%s/grouppush/%s", m.(*Config).orgName, m.(*Config).domain, d.Get("org_to_org_instance_id"), d.Id())

	org := pushgroupsreq{
		Status: d.Get("status").(string),
	}

	client := &http.Client{}

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
		return fmt.Errorf("failed to Push groups for Org2Org app for url: %s, status: %s, req: %s", url, resp.Status, bytesJson)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getPushGroupsResJson([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(s.MappingId)
	d.Set("status", s.Status)

	d.Partial(false)

	return nil
}

func resourcepushGroupsDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
