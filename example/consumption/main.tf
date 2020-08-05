provider "azurerm" {
  # whilst the `version` attribute is optional, we recommend pinning to a given version of the Provider
  version = "=2.0.0"
  features {}
}

module "actionGroup" {
  source = "github.com/AlaskaAirlines/tfmodule-azure-actiongroup.git?ref=v1.0.1"

  resource-group-name = var.resource-group-name
  appName             = "emailSample"
  environment         = "test"
  shortName           = "blah"
  enableEmail         = true
  emailName           = "TestName"
  emailAddress        = "test@alaskaair.com"
}

module "consumptionModule" {
  source = "../../."

  resource-group-name   = var.resource-group-name
  appName               = var.appName
  environment           = var.environment
  location              = var.location
  kind                  = var.kind
  actionGroupId         = module.actionGroup.action_group_id
  autoScaleNotifyEmails = var.autoScaleNotifyEmails
}
