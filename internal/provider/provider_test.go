package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider("1").InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func init() {
	testAccProvider = Provider("1")
	testAccProviders = map[string]*schema.Provider{
		"sparkpost": testAccProvider,
	}
}
