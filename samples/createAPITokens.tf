
provider "okadmin" {
        org_name = "terraform-p2"
        base_url  = "oktapreview.com"
        api_token = "123455"
}

resource "okadmin_createApiTokens" "createTokens" {
        name="test1234"
}
