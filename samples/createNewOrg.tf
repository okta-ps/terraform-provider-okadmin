provider "okadmin" {
        org_name = "terraform-p2-template"
        base_url  = "oktapreview.com"
        api_token = "123455"
	alias = "templateOrg"
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
