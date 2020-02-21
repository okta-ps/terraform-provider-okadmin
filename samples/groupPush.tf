provider "okadmin" {
        org_name = "terraform-p2"
        base_url  = "oktapreview.com"
        api_token = "123455"
}


/* user_group_id is the group ID that has to be pushed
org_to_org_instance_id is the ID of the Org2Org app created*/

resource "okadmin_pushGroups" "pushAllGroups" {
      status = "ACTIVE"
      user_group_id="00go70s9ywuTmw9fV0h7"
      org_to_org_instance_id="0oapmxfaf3vVrMzsG0h7"
}
