package provider

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/services/preview/managementpartner/mgmt/2018-02-01/managementpartner"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jongio/azidext/go/azidext"
	"go.uber.org/multierr"
)

func resourceManagementPartner() *schema.Resource {
	return &schema.Resource{
		Description:   "Configures Azure Admin Partner Link",
		CreateContext: resourceManagementPartnerCreate,
		ReadContext:   resourceManagementPartnerRead,
		UpdateContext: resourceManagementPartnerUpdate,
		DeleteContext: resourceManagementPartnerDelete,
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
			"overwrite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Overwrite existing PAL",
			},
		},
	}
}

func resourceManagementPartnerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	partnerID := d.Get("partner_id").(string)
	overwrite := d.Get("overwrite").(bool)

	mpClient, err := setupClient(ctx, tenantID, clientID, clientSecret)
	if err != nil {
		return diag.FromErr(err)
	}

	createErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := mpClient.Create(ctx, partnerID)
		// The request needs to be retried as sometimes the client secret takes time to become
		// valid even though a token is returned.
		if err != nil && strings.Contains(err.Error(), "AADSTS7000215") {
			err := fmt.Errorf("client secret is yet to be propogated (AADSTS7000215): %v", err)
			log.Printf("[DEBUG] %v", err)
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil && !overwrite {
		return diag.FromErr(fmt.Errorf("could not create management partner: %w", err))
	}

	if createErr != nil && overwrite {
		if _, err := mpClient.Update(ctx, partnerID); err != nil {
			return diag.Errorf("could not create/update management partner: %v", multierr.Combine(createErr, err))
		}
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := mpClient.Get(ctx, partnerID)
		if err != nil {
			err = fmt.Errorf("could not get management partner: %v", err)
			log.Printf("[DEBUG] %v", err)
			return resource.RetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("could not get created management partner: %v", err)
	}

	d.SetId(fmt.Sprintf("%s-%s", clientID, partnerID))
	return resourceManagementPartnerRead(ctx, d, m)
}

func resourceManagementPartnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	partnerID := d.Get("partner_id").(string)

	mpClient, err := setupClient(ctx, tenantID, clientID, clientSecret)
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := mpClient.Get(ctx, partnerID)
		// The request needs to be retried as sometimes the client secret takes time to become
		// valid even though a token is returned.
		if err != nil && strings.Contains(err.Error(), "AADSTS7000215") {
			err := fmt.Errorf("client secret is yet to be propogated (AADSTS7000215): %v", err)
			log.Printf("[DEBUG] %v", err)
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not read management partner: %w", err))
	}

	d.SetId(fmt.Sprintf("%s-%s", clientID, partnerID))
	d.Set("partner_id", partnerID)
	return nil
}

func resourceManagementPartnerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	partnerID := d.Get("partner_id").(string)

	mpClient, err := setupClient(ctx, tenantID, clientID, clientSecret)
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := mpClient.Update(ctx, partnerID)
		// The request needs to be retried as sometimes the client secret takes time to become
		// valid even though a token is returned.
		if err != nil && strings.Contains(err.Error(), "AADSTS7000215") {
			err := fmt.Errorf("client secret is yet to be propogated (AADSTS7000215): %v", err)
			log.Printf("[DEBUG] %v", err)
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not update management partner: %w", err))
	}

	return resourceManagementPartnerRead(ctx, d, m)
}

func resourceManagementPartnerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	partnerID := d.Get("partner_id").(string)

	mpClient, err := setupClient(ctx, tenantID, clientID, clientSecret)
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := mpClient.Delete(ctx, partnerID)
		// The request needs to be retried as sometimes the client secret takes time to become
		// valid even though a token is returned.
		if err != nil && strings.Contains(err.Error(), "AADSTS7000215") {
			err := fmt.Errorf("client secret is yet to be propogated (AADSTS7000215): %v", err)
			log.Printf("[DEBUG] %v", err)
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete management partner: %w", err))
	}

	d.SetId("")
	return nil
}

func setupClient(ctx context.Context, tenantID, clientID, clientSecret string) (*managementpartner.PartnerClient, error) {
	defaultScope := []string{"https://management.azure.com/.default"}

	// Configure client secret credenitals
	// Zero value options will be defaulted
	opts := &azidentity.ClientSecretCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				// A value less than zero means one try and no retries
				MaxRetries: -1,
			},
		},
	}
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, opts)
	if err != nil {
		return nil, fmt.Errorf("invalid client credentials: %w", err)
	}

	// Wait for Service Account credentials to be valid as it may take a while if just created
	timeout := 5 * time.Minute
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		opts := policy.TokenRequestOptions{
			Scopes: defaultScope,
		}
		if _, err := cred.GetToken(ctx, opts); err != nil {
			err = fmt.Errorf("could not get valid token: %w", err)
			log.Printf("[DEBUG] %v", err)
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not verify client credentials: %w", err)
	}

	// Create and return partner client
	mpClient := managementpartner.NewPartnerClient()
	mpClient.Authorizer = azidext.NewTokenCredentialAdapter(cred, defaultScope)
	mpClient.RetryAttempts = 0
	return &mpClient, nil
}
