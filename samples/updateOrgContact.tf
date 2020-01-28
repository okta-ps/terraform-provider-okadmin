 
 
 provider "okadmin" {
        org_name = "terraform-p2"
        base_url  = "oktapreview.com"
        api_token = "123455"
}

###State should be two Letters (upperCase), Country Code must be included in upper case, No trailing or leading spaces must be present

 resource "okadmin_updateOrgContactAddress" "orgContacts" {
        street_address_1="1st Ave"
        street_address_2="1st Ave"
        city= "Sanfransisco"
        state= "CA"
        zip_code="91303"
        country= "USA"
        mobile_phone = "9376546869"
        technical_phone_number = "9376546869"
}
