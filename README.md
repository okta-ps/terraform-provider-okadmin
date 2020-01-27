# Terraform Admin Provider Okta 

## Maintainers

This provider plugin is maintained by the Praveen Atluri. 

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

## Demo



## Usage

This plugin requires three inputs to run: the okta organization name, okta base url and the okta api token. 

You can specify the inputs in your tf plan:

```
provider "okadmin" {
  org_name  = <okta instance name, e.g. dev-XXXXXX>
  api_token = <okta instance api token with the Administrator role>
  base_url  = <okta base url, e.g. oktapreview.com>
}
```

## Examples

As we build out resources we build concomitant acceptance tests that require use to create resource config that actually creates and modifies real resources. We decided to put these test fixtures to good use and provide them [as examples here.](./examples)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/patluri43/terraform-provider-okadmin`

cd $GOPATH/src/github.com/patluri43/terraform-provider-okadmin
./build.sh


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
    website= "https://google456.com"
    secondary_email= false
    secondary_image= true
    app_error_redirect_url = "https://google.com"
    portal_error_redirect_url = "https://google.com"
    portal_signout_url = "https://google.com"
    interstitial_min_wait_time = 900
}
```

## Disclaimer

There are particular resources and settings that are not exposed on Okta's public API. Please submit an issue if you find one not listed here.

### Org Settings

- Org level customization settings.


## Common Errors


## Developing the Provider



### Best Practices
