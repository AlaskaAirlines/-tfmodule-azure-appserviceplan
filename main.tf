provider "azurerm" {
  # whilst the `version` attribute is optional, we recommend pinning to a given version of the Provider
  version = "=2.0.0"
  features {}
}

locals {
  baseName = "${var.appName}-${var.environment}-sharedplan"
}

data "azurerm_resource_group" "rg" {
  name = var.resource-group-name
}

resource "azurerm_app_service_plan" "sharedplan" {
  count = var.instanceCount

  name                = "${local.baseName}-${count.index}-${var.location}"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = var.location

  tags     = data.azurerm_resource_group.rg.tags
  kind     = var.kind
  reserved = var.kind == "Linux" || var.kind == "linux" ? true : var.kind == "Windows" || var.kind == "windows" || var.kind == "App" || var.kind == "app" ? false : var.reserved

  dynamic "sku" {
    for_each = var.kind == "FunctionApp" ? ["sku"] : []
    content {
      tier = "Dynamic"
      size = "Y1"
    }
  }

  dynamic "sku" {
    for_each = var.kind != "FunctionApp" ? ["sku"] : []
    content {
      tier     = var.skuTier
      size     = var.skuSize
      capacity = var.skuCapacity
    }
  }
}

# resource "azurerm_monitor_metric_alert" "planAlert" {
#   count = var.instanceCount

#   name                = "${local.baseName}-${count.index}git -alerts"
#   resource_group_name = data.azurerm_resource_group.rg.name
#   scopes              = var.appServicePlanIds
#   description         = "Metric alerts for an app service plan"

#   criteria {
#     metric_namespace = "Microsoft.Web/serverfarms"
#     metric_name      = "CpuPercentage"
#     aggregation      = "Average"
#     operator         = "GreaterThan"
#     threshold        = var.cpuThreshold
#   }

#   criteria {
#     metric_namespace = "Microsoft.Web/serverfarms"
#     metric_name      = "DiskQueueLength"
#     aggregation      = "Average"
#     operator         = "GreaterThan"
#     threshold        = 100
#   }

#   criteria {
#     metric_namespace = "Microsoft.Web/serverfarms"
#     metric_name      = "MemoryPercentage"
#     aggregation      = "Average"
#     operator         = "GreaterThan"
#     threshold        = 90
#   }

#   criteria {
#     metric_namespace = "Microsoft.Web/serverfarms"
#     metric_name      = "HttpQueueLength"
#     aggregation      = "Average"
#     operator         = "GreaterThan"
#     threshold        = 100
#   }

#   action {
#     action_group_id = var.actionGroupId
#   }
# }

# resource "" "autoscale" {
#   count = var.kind == "FunctionApp" || var.kind == "elastic" ? 0 : var.count
# }