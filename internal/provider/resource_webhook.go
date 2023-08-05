package provider

import (
	"context"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateCreate,
		ReadContext:   dataSourceWebhookRead,
		UpdateContext: resourceTemplateUpdate,
		DeleteContext: resourceTemplateDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Unique alphanumeric ID used to reference the webhook.",
				Type:        schema.TypeString,
				ConfigMode:  true,
			},
			"name": {
				Description: "Editable display name. At a minimum, id or name is required upon creation. Does not have to be unique. Maximum length - 1024 bytes",
				Type:        schema.TypeString,
				Required:    true,
			},
			"target": {
				Description: "URL of the target to which to POST event batches. Only ports 80 for http and 443 for https can be set.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"events": {
				Description: "",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					MinItems: 1,
				},
			},
			"active": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"auth_type": {
				Description: "",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			// For simplicity, the provider can only import published templates
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resouceWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// get the sparkpost client
	client := m.(*sp.Client)

	webhookID := d.Id()

	// SparkPost requires that we request a draft or published version of the template
	draft, ok := d.GetOk("draft")
	if !ok {
		draft = false
	}

	hook := &sp.WebhookDetailWrapper{
		ID: webhookID,
	}

	_, err := client.WebhookDetailContext(ctx, hook)
	if err != nil {
		return diag.FromErr(err)
	}

	err = setWebhookResourceData(d, template)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(webhookID)

	return diags
}

func buildWebhook(webhookID string, d *schema.ResourceData) *sp.WebhookItem {

	hook := &sp.WebhookItem{
		ID:       webhookID,
		Name:     d.Get("name").(string),
		Target:   d.Get("target").(string),
		Events:   d.Get("events").([]string),
		AuthType: d.Get("auth_type").(string),
	}

	return hook
}

func publishTemplate(ctx context.Context, d *schema.ResourceData, client *sp.Client, templateID string) error {
	// automatically publish the template
	_, err := client.WebTemplatePublishContext(ctx, templateID)
	if err != nil {
		return err
	}

	return nil
}

func resourceWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sp.Client)

	templateID := getTemplateID(d)

	template := buildTemplate(templateID, d)

	id, _, err := client.TemplateCreateContext(ctx, template)
	if err != nil {
		return diag.FromErr(err)
	}

	if template.Published {
		err = publishTemplate(ctx, d, client, templateID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Ensure the read looks for the correct copy of the template
	d.Set("draft", !template.Published)

	d.SetId(id)

	return resourceTemplateRead(ctx, d, m)
}

func resourceTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*sp.Client)

	templateID := d.Id()

	published := d.Get("published").(bool)

	template := buildTemplate(templateID, d)

	var publishUpdate, updatePublished bool

	if d.HasChange("published") && published {
		// was draft now published
		publishUpdate = true

		// WARNING: it's undocumented, but the ?update_published param on the PUT can be overridden by
		// a `published` field in the body. Since we want to update the draft here, we MUST set the
		// published field to false in the PUT body.
		// https://developers.sparkpost.com/api/templates/#templates-put-update-a-published-template
		template.Published = false
	} else if published {
		// was published, no change
		updatePublished = true
	}

	// Update the template. `updatePublished` controls whether to update the published or draft copy.
	_, err := client.TemplateUpdateContext(ctx, template, updatePublished)
	if err != nil {
		return diag.FromErr(err)
	}

	// Publish it if we're going from draft -> published
	if publishUpdate {
		err = publishTemplate(ctx, d, client, templateID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Ensure the read looks for the correct copy of the template
	d.Set("draft", !published)

	return resourceTemplateRead(ctx, d, m)
}

func resourceTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*sp.Client)

	templateID := d.Id()

	_, err := client.TemplateDeleteContext(ctx, templateID)
	if err != nil {
		return diag.FromErr(err)
	}

	// mark as deleted
	d.SetId("")

	return diags
}

func getTemplateID(d *schema.ResourceData) string {
	templateID, ok := d.GetOk("template_id")
	if !ok {
		return ""
	}
	return templateID.(string)
}
