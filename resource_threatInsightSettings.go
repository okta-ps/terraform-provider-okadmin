package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceThreatInsightSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceThreatInsightSettingsCreate,
		Read:   resourceThreatInsightSettingsRead,
		Update: resourceThreatInsightSettingsUpdate,
		Delete: resourceThreatInsightSettingsDelete,

		Schema: map[string]*schema.Schema{
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"exclude_zones": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

type threatinsight struct {
	Action        string   `json:"action,omitempty"`
	Exclude_zones []string `json:"excludeZones,omitempty"`
}

func resourceThreatInsightSettingsCreate(d *schema.ResourceData, m interface{}) error {
	return resourceThreatInsightSettingsUpdate(d, m)
}

func getThreatInsightsJson(body []byte) (*threatinsight, error) {
	var s = new(threatinsight)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error Unmarshalling Json Response:", err)
	}
	return s, err
}

func resourceThreatInsightSettingsRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("https://%s.%s/api/v1/threats/configuration", m.(*Config).orgName, m.(*Config).domain)

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
		return fmt.Errorf("failed to get Threat Insight Settings for url: %s, status: %s", url, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getThreatInsightsJson([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.SetId(fmt.Sprintf("%s.%s-ThreatInsights", m.(*Config).orgName, m.(*Config).domain))
	d.Set("action", s.Action)
	d.Set("exclude_zones", s.Exclude_zones)
	return nil
}

func resourceThreatInsightSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	url := fmt.Sprintf("https://%s.%s/api/v1/threats/configuration", m.(*Config).orgName, m.(*Config).domain)

	client := &http.Client{}

	org := threatinsight{
		Action:        d.Get("action").(string),
		Exclude_zones: convertInterfaceToStringSet(d.Get("exclude_zones")),
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
		return fmt.Errorf("failed to Modify Threat Insight Page Settings for url: %s, status: %s, req: %s", url, resp.Status, bytesJson)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getThreatInsightsJson([]byte(body))
	if err == nil {
		fmt.Println("Error getting the Json Response:", s)
	}

	d.Partial(false)

	return resourceThreatInsightSettingsRead(d, m)
}

func resourceThreatInsightSettingsDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func convertInterfaceToStringSet(purportedSet interface{}) []string {
	return convertInterfaceToStringArr(purportedSet.(*schema.Set).List())
}

func convertInterfaceToStringArr(purportedList interface{}) []string {
	var arr []string
	rawArr, ok := purportedList.([]interface{})

	if ok {
		arr = make([]string, len(rawArr))
		for i, thing := range rawArr {
			arr[i] = thing.(string)
		}
	}

	return arr
}
