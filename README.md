# Terrraform Provider PAL
Provider to configure [Partner Admin Link](https://docs.microsoft.com/en-us/azure/cost-management-billing/manage/link-partner-id) for indiviual Service Principal.

## How To
The provider provides the resource `pal_management_partner` which creates a PAL binding between
a Service Principal and a parnter id. The provider does not take any parameters, instead the
id and credentials for the Service Principal are given for each resource. This is required as
each Service Principal needs to configure the PAL individually.

```hcl
provider "azurerm" {
  version = "=2.14"
  features {}
}

provider "azuread" {
  version = "=0.10.0"
}

resource "azuread_application" "aadApp" {
  name = "example"
}

resource "azuread_service_principal" "aadSp" {
  application_id = azuread_application.aadApp.application_id
}

resource "random_password" "aadSpSecret" {
  length           = 24
  special          = true
  override_special = "!-_="

  keepers = {
    service_principal = azuread_service_principal.aadSp.id
  }
}

resource "azuread_application_password" "aadSpSecret" {
  application_object_id = azuread_application.aadApp.id
  value                 = random_password.aadSpSecret.result
  end_date              = timeadd(timestamp(), "87600h") # 10 years

  lifecycle {
    ignore_changes = [
      end_date
    ]
  }
}

provider "pal" {}

data "azurerm_client_config" "current" {}

resource "pal_management_partner" "foobar" {
  tenant_id = data.azurerm_client_config.current.tenant_id
  client_id = azuread_service_principal.aadSp.application_id
  client_secret = random_password.aadSpSecret.result
  partner_id = "2501985"
}
```

## Issues
A issue is that the provider can be experienced as slightly slower than others. This will
mostly be due to the fact that each resource needs to setup its own authentication client.
Currently this is the only solution to the problem as a provider cant be iterated over,
so each resource needs to be given authentication credentials instead of providig them
to the provider.
