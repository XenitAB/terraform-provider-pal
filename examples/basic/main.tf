terraform {
  required_version = ">=0.13.5"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">=2.35.0"
    }
    azuread = {
      source = "hashicorp/azuread"
      version = ">=1.1.1"
    }
    pal = {
      source  = "xenitab/pal"
      version = "0.0.0-dev"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "random_password" "this" {
  length           = 48
  special          = true
  override_special = "!-_="
}

resource "azuread_application" "this" {
  name = "terraform-provider-pal-test-app"
}

resource "azuread_application_password" "this" {
  application_object_id = azuread_application.this.id
  value                 = random_password.this.result
  end_date              = timeadd(timestamp(), "87600h") # 10 years

  lifecycle {
    ignore_changes = [
      end_date
    ]
  }
}

data "azurerm_subscription" "current" {}

resource "azuread_service_principal" "this" {
  application_id = azuread_application.this.application_id
}

resource "pal_management_partner" "this" {
  tenant_id     = data.azurerm_subscription.current.tenant_id
  client_id     = azuread_service_principal.this.application_id
  client_secret = random_password.this.result
  partner_id    = var.partner_id
}
