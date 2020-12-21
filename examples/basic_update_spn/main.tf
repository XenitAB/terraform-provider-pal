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

data "azurerm_subscription" "current" {}


data "azuread_service_principal" "this" {
  display_name = var.service_principal_name
}

data "azuread_application" "this" {
  name = var.service_principal_name
}

resource "random_password" "this" {
  length           = 48
  special          = true
  override_special = "!-_="

  keepers = {
    service_principal = data.azuread_service_principal.this.id
  }
}

resource "azuread_application_password" "this" {
  application_object_id = data.azuread_application.this.object_id
  value                 = random_password.this.result
  end_date              = timeadd(timestamp(), "87600h") # 10 years

  lifecycle {
    ignore_changes = [
      end_date
    ]
  }
}

resource "pal_management_partner" "this" {
  depends_on = [azuread_application_password.this]

  tenant_id     = data.azurerm_subscription.current.tenant_id
  client_id     = data.azuread_service_principal.this.application_id
  client_secret = random_password.this.result
  partner_id    = var.partner_id
  overwrite     = true
}
