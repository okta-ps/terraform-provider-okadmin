package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceInterstitialSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceInterstitialSettingsCreate,
		Read:   resourceInterstitialSettingsRead,
		Update: resourceInterstitialSettingsUpdate,
		Delete: resourceInterstitialSettingsDelete,

		Schema: map[string]*schema.Schema{
			"interstitial_page_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"interstitial_page_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceInterstitialSettingsCreate(d *schema.ResourceData, m interface{}) error {
	oktaInterstitialPageId := d.Get("interstitial_page_id").(string)
	d.SetId(oktaInterstitialPageId)
	return resourceInterstitialSettingsRead(d, m)
}

func resourceInterstitialSettingsRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceInterstitialSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceInterstitialSettingsRead(d, m)
}

func resourceInterstitialSettingsDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
