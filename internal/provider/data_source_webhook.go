package provider

import (
	"context"
	"fmt"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWebhook() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebhookRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Unique alphanumeric ID used to reference the webhook.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Editable display name. At a minimum, id or name is required upon creation. Does not have to be unique. Maximum length - 1024 bytes",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"target": {
				Description: "URL of the target to which to POST event batches. Only ports 80 for http and 443 for https can be set.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"events": {
				Description: "",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"active": {
				Description: "",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"exception_subaccounts": {
				Description: "",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"auth_type": {
				Description: "",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*sp.Client)

	webhookID := d.Get("id").(string)

	webhook := &sp.WebhookDetailWrapper{
		ID: webhookID,
	}

	_, err := client.WebhookDetailContext(ctx, webhook)

	fmt.Printf("%+v\n", webhook.Results)

	if err != nil {
		return diag.FromErr(err)
	}

	err = setWebhookResourceData(d, webhook.Results)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(webhookID)

	return diags
}

func setWebhookResourceData(d *schema.ResourceData, webhook *sp.WebhookItem) error {

	d.Set("id", webhook.ID)
	d.Set("name", webhook.Name)
	d.Set("target", webhook.Target)
	d.Set("events", webhook.Events)
	d.Set("auth_type", webhook.AuthType)
	d.Set("events", webhook.Events)

	return nil
}
