output "sharedplanids" {
  value = "${azurerm_app_service_plan.sharedplan.*.id}"
}