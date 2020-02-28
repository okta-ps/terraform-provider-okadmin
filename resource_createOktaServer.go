package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCreateOktaServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateOktaServerCreate,
		Read:   resourceCreateOktaServerRead,
		Update: resourceCreateOktaServerUpdate,
		Delete: resourceCreateOktaServerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// This allows importing an initial token for an org already
				// created via external method.  Subsequent run will "update"
				// state to match config values, but since Update function does
				// nothing it will have no effect.
				if d.Id() != "env" {
					return nil, errors.New("Only \"env\" import type is supported")
				}
				d.Set("token", os.Getenv("IMPORT_ORG_token"))
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
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
				Required: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"login": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mobile_phone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"recovery_question": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"recovery_answer": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"secondary_email": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"secondary_image": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"app_error_redirect_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"portal_error_redirect_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"portal_signout_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"interstitial_min_wait_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"token": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"token_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"edition_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type Profile struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Email       string `json:"email,omitempty"`
	Login       string `json:"login,omitempty"`
	MobilePhone string `json:"mobilePhone,omitempty"`
}
type Password struct {
	Value string `json:"value,omitempty"`
}
type RecoveryQuestion struct {
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
}

type Credentials struct {
	Password         *Password         `json:"password,omitempty"`
	RecoveryQuestion *RecoveryQuestion `json:"recovery_question,omitempty"`
}

type Admin struct {
	Profile     *Profile     `json:"profile,omitempty"`
	Credentials *Credentials `json:"credentials,omitempty"`
}

type App struct {
	ErrorRedirectURL        interface{} `json:"errorRedirectUrl,omitempty"`
	InterstitialMinWaitTime int         `json:"interstitialMinWaitTime,omitempty"`
}
type Attributes struct {
	SecondaryEmail interface{} `json:"secondaryEmail,omitempty"`
	SecondaryImage interface{} `json:"secondaryImage,omitempty"`
}
type UserAccount struct {
	Attributes *Attributes `json:"attributes,omitempty"`
}
type Portal struct {
	ErrorRedirectURL interface{} `json:"errorRedirectUrl,omitempty"`
	SignOutURL       interface{} `json:"signOutUrl,omitempty"`
}
type Settings struct {
	App         *App         `json:"app,omitempty"`
	UserAccount *UserAccount `json:"userAccount,omitempty"`
	Portal      *Portal      `json:"portal,omitempty"`
}

type orgsettings struct {
	ID        string    `json:"id,omitempty"`
	Subdomain string    `json:"subdomain"`
	Name      string    `json:"name"`
	Website   string    `json:"website"`
	EditionId string    `json:"editionId"`
	Settings  *Settings `json:"settings,omitempty"`
	Token     string    `json:"token,omitempty"`
	TokenType string    `json:"tokenType,omitempty"`
	Admin     *Admin    `json:"admin,omitempty"`
}

func resourceCreateOktaServerCreate(d *schema.ResourceData, m interface{}) error {

	url := fmt.Sprintf("https://%s.%s/api/v1/orgs", m.(*Config).orgName, m.(*Config).domain)

	client := &http.Client{}

	org := &orgsettings{
		Subdomain: d.Get("subdomain").(string),
		Name:      d.Get("name").(string),
		Website:   d.Get("website").(string),
		EditionId: d.Get("edition_id").(string),
		Admin: &Admin{
			Profile: &Profile{
				FirstName:   d.Get("first_name").(string),
				LastName:    d.Get("last_name").(string),
				Email:       d.Get("email").(string),
				Login:       d.Get("login").(string),
				MobilePhone: d.Get("mobile_phone").(string),
			},
			Credentials: &Credentials{
				Password: &Password{
					Value: d.Get("password").(string),
				},
				RecoveryQuestion: &RecoveryQuestion{
					Question: d.Get("recovery_question").(string),
					Answer:   d.Get("recovery_answer").(string),
				},
			},
		},
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
	} else if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to update Org Details for url: %s, status: %s, req: %s", url, resp.Status, bytesJson)
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
	d.Set("website", s.Website)
	d.Set("name", s.Name)
	d.Set("subdomain", s.Subdomain)
	d.Set("secondary_email", s.Settings.UserAccount.Attributes.SecondaryEmail)
	d.Set("secondary_image", s.Settings.UserAccount.Attributes.SecondaryImage)
	d.Set("app_error_redirect_url", s.Settings.App.ErrorRedirectURL)
	d.Set("portal_error_redirect_url", s.Settings.Portal.ErrorRedirectURL)
	d.Set("portal_signout_url", s.Settings.Portal.SignOutURL)
	d.Set("interstitial_min_wait_time", s.Settings.App.InterstitialMinWaitTime)
	d.Set("token", s.Token)
	d.Set("token_type", s.TokenType)

	return nil
}

func resourceCreateOktaServerRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceCreateOktaServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceCreateOktaServerRead(d, m)
}

func resourceCreateOktaServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
