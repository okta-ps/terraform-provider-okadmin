

provider "okadmin" {
        org_name = "terraform-p2"
        base_url  = "oktapreview.com"
        api_token = "123455"
}

resource "okadmin_server" "MyOrg" {
    subdomain= "terraform-p2"
    name= "terraform-p2s"
    website= "https://client1.greenlake.com"
    secondary_email= true
    secondary_image= true
    error_redirect_url = "https://client.greenlake.com"
    interstitial_min_wait_time = 1200 
}

