
provider "okadmin" {
        org_name = "terraform-p2"
        base_url  = "oktapreview.com"
        api_token = "123455"
}

/* action can be none/audit/block */
resource "okadmin_updateTheatInsightSettings" "tiSettings" {
  action= "block"
  exclude_zones =  [ ]
}
