module "basicModule" {
  source = "../../."

  resource-group-name = var.resource-group-name
  appName             = var.appName
  environment         = var.environment
  location            = var.location
}
