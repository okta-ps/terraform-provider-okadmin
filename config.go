package main

import (
	"context"
	"fmt"
	"net/http"

	// "github.com/articulate/terraform-provider-okta/sdk"

	// articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/okta/okta-sdk-golang/okta"
	// "github.com/okta/okta-sdk-golang/okta"
)

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "Okta Org Creator Terraform Provider")
	return adt.T.RoundTrip(req)
}

type (
	// AddHeaderTransport used to tack on default headers to outgoing requests
	AddHeaderTransport struct {
		T http.RoundTripper
	}

	// Config contains our provider schema values and Okta clients
	Config struct {
		orgName  string
		domain   string
		apiToken string
		// retryCount   int
		// parallelism  int
		// waitForReset bool
		// backoff      bool
		// minWait      int
		// maxWait      int

		supplementClient *ApiSupplement
	}
)

func (c *Config) loadAndValidate() error {
	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport = logging.NewTransport("Okta", httpClient.Transport)

	orgUrl := fmt.Sprintf("https://%v.%v", c.orgName, c.domain)

	client, err := okta.NewClient(
		context.Background(),
		okta.WithOrgUrl(orgUrl),
		okta.WithToken(c.apiToken),
		okta.WithCache(false),
		// okta.WithBackoff(c.backoff),
		// okta.WithMinWait(time.Duration(c.minWait)*time.Second),
		// okta.WithMaxWait(time.Duration(c.maxWait)*time.Second),
		// okta.WithRetries(int32(c.retryCount)),
		okta.WithHttpClient(*httpClient),
	)
	if err != nil {
		return err
	}

	c.supplementClient = &ApiSupplement{
		BaseURL:         orgUrl,
		Client:          httpClient,
		Token:           c.apiToken,
		RequestExecutor: client.GetRequestExecutor(),
	}

	return nil
}
