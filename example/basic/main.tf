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

module "basicModule" {
  source = "../../."

  resource-group-name = var.resource-group-name
  appName             = var.appName
  environment         = var.environment
  location            = var.location
  actionGroupId       = module.actionGroup.action_group_id
}
