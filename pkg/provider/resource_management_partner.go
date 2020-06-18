package provider

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/managementpartner/mgmt/2018-02-01/managementpartner"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceManagementPartner() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagementPartnerCreate,
		Read:   resourceManagementPartnerRead,
		Update: resourceManagementPartnerUpdate,
		Delete: resourceManagementPartnerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceManagementPartnerImport,
		},

		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Tenant ID which should be used.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Client ID which should be used.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},
			"partner_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the partner to link to.",
			},
		},
	}
}

func resourceManagementPartnerCreate(d *schema.ResourceData, m interface{}) error {
	clientID := d.Get("client_id").(string)
	partnerID := d.Get("partner_id").(string)

	mpClient, err := setupClient(d)
	if err != nil {
		return err
	}

	// Needed as the SA password needs some time to become valid
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		if err := resourceManagementPartnerRead(d, m); err != nil {
			if _, err := mpClient.Create(context.Background(), partnerID); err != nil {
				return resource.RetryableError(err)
			}
		} else {
			if _, err := mpClient.Update(context.Background(), partnerID); err != nil {
				return resource.RetryableError(err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(clientID + "-" + partnerID)
	return resourceManagementPartnerRead(d, m)
}

func resourceManagementPartnerRead(d *schema.ResourceData, m interface{}) error {
	clientID := d.Get("client_id").(string)
	partnerID := d.Get("partner_id").(string)

	mpClient, err := setupClient(d)
	if err != nil {
		return err
	}

	if _, err := mpClient.Get(context.Background(), partnerID); err != nil {
		return err
	}

	d.SetId(clientID + "-" + partnerID)
	d.Set("partner_id", partnerID)
	return nil
}

func resourceManagementPartnerUpdate(d *schema.ResourceData, m interface{}) error {
	partnerID := d.Get("partner_id").(string)

	mpClient, err := setupClient(d)
	if err != nil {
		return err
	}

	if _, err := mpClient.Update(context.Background(), partnerID); err != nil {
		return err
	}

	return resourceManagementPartnerRead(d, m)
}

func resourceManagementPartnerDelete(d *schema.ResourceData, m interface{}) error {
	partnerID := d.Get("partner_id").(string)

	mpClient, err := setupClient(d)
	if err != nil {
		return err
	}

	if _, err := mpClient.Delete(context.Background(), partnerID); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceManagementPartnerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceManagementPartnerRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func setupClient(d *schema.ResourceData) (*managementpartner.PartnerClient, error) {
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	mpClient := managementpartner.NewPartnerClient()
	authorizer, err := auth.NewClientCredentialsConfig(clientID, clientSecret, tenantID).Authorizer()
	if err != nil {
		return nil, err
	}
	mpClient.Authorizer = authorizer

	return &mpClient, nil
}
