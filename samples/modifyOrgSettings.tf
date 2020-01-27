
provider "okadmin" {
        org_name = "terraform-p2"
        base_url  = "oktapreview.com"
        api_token = "123455"
}

resource "okadmin_updateServerSettings" "MyOrg" {
    name= "terraform-p2s"
    website= "https://okta.com"
    secondary_email= false
    secondary_image= false
    app_error_redirect_url = "https://okta.com"
    portal_error_redirect_url = "https://okta.com"
    portal_signout_url = "https://okta.com"
    interstitial_min_wait_time = 1200 
}