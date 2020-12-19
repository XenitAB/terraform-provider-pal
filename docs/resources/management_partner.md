---
page_title: "pal_management_partner Resource - terraform-provider-pal"
subcategory: ""
description: |-
  Configures Azure Admin Partner Link
---

# Resource `pal_management_partner`

Configures Azure Admin Partner Link



## Schema

### Required

- **client_id** (String) The Client ID which should be used.
- **client_secret** (String, Sensitive) The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.
- **partner_id** (String) The ID of the partner to link to.
- **tenant_id** (String) The Tenant ID which should be used.

### Optional

- **id** (String) The ID of this resource.
- **overwrite** (Boolean) Overwrite existing PAL Defaults to `false`.


