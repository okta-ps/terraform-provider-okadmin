# Terraform Okta Admin Provider

## Maintainers



## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

## Usage

This plugin requires three inputs to run: the okta organization name, okta base url and the okta api token. 

You can specify the inputs in your tf plan:

```
provider "okadmin" {
  org_name  = <e.g. subdomain (dev-xxxxx)>
  api_token = <Admin apikey for the okta instance>
  base_url  = <e.g. oktapreview.com (Preview) or okta.com (PROD)>
}
```


## Building The Provider

Clone repository to: `$GOPATH/src/github.com/patluri43/terraform-provider-okadmin`

cd $GOPATH/src/github.com/patluri43/terraform-provider-okadmin
```
./build.sh
```

## Using the provider

Example terraform plan:

```
provider "okadmin" {
        org_name = "terraform-p2"
        base_url  = "oktapreview.com"
        api_token = "00jn1UE7gvcY92ohwApoUbRi7GsV93GRQRuZe8W-CC"
}

resource "okadmin_updateServerSettings" "MyOrg" {
    name= "terraform-p2"
    website= "https://okta456.com"
    secondary_email= false
    secondary_image= true
    app_error_redirect_url = "https://okta.com"
    portal_error_redirect_url = "https://okta.com"
    portal_signout_url = "https://okta.com"
    interstitial_min_wait_time = 900
}
```

For More samples. Please go to 

```
https://github.com/patluri43/terraform-provider-okadmin/tree/master/samples
```

## Supported Resources by Provider

			"okadmin_updateServerSettings"
			"okadmin_interstitialSettings"
			"okadmin_createOktaServer"
			"okadmin_updateTechnicalContacts"
			"okadmin_updateTheatInsightSettings"
            "okadmin_updateOrgContactAddress"
			"okadmin_createApiTokens"


## Disclaimer

Provider is consuming Okta Private API's which are subject to change with out notice. 