

provider "okadmin" {
        org_name = "terraform-p2"
        base_url  = "oktapreview.com"
        api_token = "123455"
}


resource "okadmin_updateTechnicalContacts" "myContacts" {
   technical_contact_id= "00uosmv4muGzGXEP20h7"
   phone_number= "+12345678901"
   help_url= "https://okta123.com"

}

