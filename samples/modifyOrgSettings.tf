provider "okadmin" {
        org_name = "terraform-p2-template"
        base_url  = "oktapreview.com"
        api_token = "123455"
	alias = "templateOrg"
}


provider "okadmin" {
        org_name = "terraform-p2"
        base_url  = "oktapreview.com"
        api_token = "123455"
}

resource "okadmin_updateServerSettings" "MyOrg" {
    name= "terraform-p2s"
    website= "https://google.com"
    secondary_email= false
    secondary_image= false
    app_error_redirect_url = "https://google.com"
    portal_error_redirect_url = "https://google.com"
    portal_signout_url = "https://google.com"
    interstitial_min_wait_time = 1200 
}

resource "okadmin_interstitialSettings" "interstitialSettings" {
   interstitial_page_enabled= false
}

resource "okadmin_createOktaServer" "myOktaServer" {
    subdomain="terraform-p2"
    name= "terraform-p2s"
    website= "https://google.com"
    first_name= "user"
    last_name = "ln"
    email = "user.ln@domain.com"
    login = "user.ln@domain.com"
    mobile_phone = "+12345678901"
    password = "password"
    recovery_question =  "pet's name?"
    recovery_answer = "pet"
    provider = okadmin.templateOrg
}


resource "okadmin_updateTechnicalContacts" "myContacts" {
   technical_contact_id= "00uosmv4muGzGXEP20h7"
   phone_number= "+12345678901"
   help_url= "https://google123.com"

}

/* action can be none/audit/block */
resource "okadmin_updateTheatInsightSettings" "tiSettings" {
  action= "block"
  exclude_zones =  [ ]
}

resource "okadmin_createApiTokens" "createTokens" {
        name="test1234"
}
