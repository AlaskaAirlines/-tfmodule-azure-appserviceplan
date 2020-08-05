provider "azurerm" {
  # whilst the `version` attribute is optional, we recommend pinning to a given version of the Provider
  version = "=2.18.0"
  features {}
}

module "actionGroup" {
  source = "github.com/AlaskaAirlines/tfmodule-azure-actiongroup.git?ref=v1.0.2"

  resource-group-name = var.resource-group-name
  appName             = "emailSample"
  environment         = "test"
  shortName           = "blah"
  enableEmail         = true
  emailName           = "TestName"
  emailAddress        = "test@alaskaair.com"
}

module "basicModule" {
  source = "../../."

  resource-group-name   = var.resource-group-name
  appName               = var.appName
  environment           = var.environment
  location              = var.location
  actionGroupId         = module.actionGroup.action_group_id
  autoScaleNotifyEmails = var.autoScaleNotifyEmails
}
