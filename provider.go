package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_ORG_NAME", nil),
				Description: "The organization to manage in Okta.",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_API_TOKEN", nil),
				Description: "API Token granting privileges to Okta API.",
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_BASE_URL", "okta.com"),
				Description: "The Okta url. (Use 'oktapreview.com' for Okta testing)",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"okadmin_updateServerSettings":       resourceupdateServerSettings(),
			"okadmin_interstitialSettings":       resourceInterstitialSettings(),
			"okadmin_createOktaServer":           resourceCreateOktaServer(),
			"okadmin_updateTechnicalContacts":    resourceUpdateTechContacts(),
			"okadmin_updateTheatInsightSettings": resourceThreatInsightSettings(),
			"okadmin_createApiTokens":            resourceCreateAPITokens(),
			"okadmin_updateOrgContactAddress":    resourceupdateOrgContacts(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing Okta Org Creator client")

	config := Config{
		templateOrgName: d.Get("org_name").(string),
		domain:          d.Get("base_url").(string),
		apiToken:        d.Get("api_token").(string),
	}
	if err := config.loadAndValidate(); err != nil {
		return nil, fmt.Errorf("[ERROR] Error initializing the Okta SDK clients: %v", err)
	}
	return &config, nil
}
