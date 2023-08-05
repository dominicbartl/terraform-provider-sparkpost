package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccItem_Webhook(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWebhook(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.sparkpost_webhook.sample", "name", "Molzait DEV"),
				),
			},
		},
	})
}

func testAccCheckWebhook() string {
	return fmt.Sprintf(`
		data "sparkpost_webhook" "sample" {
  			id   = "e4093700-7733-11ed-9039-2f08fa6ecfde"
		}
	`)
}
