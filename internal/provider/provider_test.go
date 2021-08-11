package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/utilitywarehouse/terraform-provider-nebraska/nebraska"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"nebraska": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("NEBRASKA_ENDPOINT") == "" {
		t.Fatal("NEBRASKA_ENDPOINT must be set for acceptance tests")
	}

	if os.Getenv("NEBRASKA_APPLICATION_ID") == "" {
		if err := os.Setenv("NEBRASKA_APPLICATION_ID", nebraska.FlatcarApplicationID); err != nil {
			t.Fatal(err)
		}
	}
}
