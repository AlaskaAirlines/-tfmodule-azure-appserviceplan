locals {
  baseName = "${var.appName}-${var.environment}-sharedplan"
  defaultAutoscaleRules = [
    {
      metricName      = "CpuPercentage",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT5M",
      timeAggregation = "Average",
      operator        = "GreaterThan",
      threshold       = 90
      direction       = "Increase",
      type            = "ChangeCount",
      value           = 2,
      cooldown        = "PT5M"
    },
    {
      metricName      = "CpuPercentage",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT5M",
      timeAggregation = "Average",
      operator        = "GreaterThan",
      threshold       = 75
      direction       = "Increase",
      type            = "ChangeCount",
      value           = 1,
      cooldown        = "PT5M"
    },
    {
      metricName      = "CpuPercentage",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT15M",
      timeAggregation = "Average",
      operator        = "LessThanOrEqual",
      threshold       = 50
      direction       = "Decrease",
      type            = "ChangeCount",
      value           = 1,
      cooldown        = "PT15M"
    },
    {
      metricName      = "HttpQueueLength",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT5M",
      timeAggregation = "Average",
      operator        = "GreaterThan",
      threshold       = 100
      direction       = "Increase",
      type            = "ChangeCount",
      value           = 1,
      cooldown        = "PT5M"
    },
    {
      metricName      = "HttpQueueLength",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT15M",
      timeAggregation = "Average",
      operator        = "LessThanOrEqual",
      threshold       = 50
      direction       = "Decrease",
      type            = "ChangeCount",
      value           = 1,
      cooldown        = "PT15M"
    },
    {
      metricName      = "HttpQueueLength",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT5M",
      timeAggregation = "Average",
      operator        = "GreaterThan",
      threshold       = 200
      direction       = "Increase",
      type            = "ChangeCount",
      value           = 2,
      cooldown        = "PT5M"
    },
    {
      metricName      = "MemoryPercentage",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT5M",
      timeAggregation = "Average",
      operator        = "GreaterThan",
      threshold       = 85
      direction       = "Increase",
      type            = "ChangeCount",
      value           = 1,
      cooldown        = "PT5M"
    },
    {
      metricName      = "MemoryPercentage",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT15M",
      timeAggregation = "Average",
      operator        = "LessThanOrEqual",
      threshold       = 65
      direction       = "Decrease",
      type            = "ChangeCount",
      value           = 1,
      cooldown        = "PT15M"
    },
    {
      metricName      = "DiskQueueLength",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT5M",
      timeAggregation = "Average",
      operator        = "GreaterThan",
      threshold       = 100
      direction       = "Increase",
      type            = "ChangeCount",
      value           = 1,
      cooldown        = "PT5M"
    },
    {
      metricName      = "DiskQueueLength",
      timeGrain       = "PT1M",
      statistic       = "Average",
      timeWindow      = "PT15M",
      timeAggregation = "Average",
      operator        = "LessThanOrEqual",
      threshold       = 50
      direction       = "Decrease",
      type            = "ChangeCount",
      value           = 1,
      cooldown        = "PT15M"
    }
  ]
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
  reserved = lower(var.kind) == "linux" ? true : lower(var.kind) == "windows" || lower(var.kind) == "app" ? false : var.reserved

  dynamic "sku" {
    for_each = lower(var.kind) == "functionapp" ? ["sku"] : []
    content {
      tier = "Dynamic"
      size = "Y1"
    }
  }

  dynamic "sku" {
    for_each = lower(var.kind) != "functionapp" ? ["sku"] : []
    content {
      tier     = var.skuTier
      size     = var.skuSize
      capacity = var.skuCapacity
    }
  }
}

resource "azurerm_monitor_metric_alert" "planAlert" {
  name                = "${local.baseName}-alerts"
  resource_group_name = data.azurerm_resource_group.rg.name
  scopes              = azurerm_app_service_plan.sharedplan.*.id
  description         = "Metric alerts for an app service plan"

  criteria {
    metric_namespace = "Microsoft.Web/serverfarms"
    metric_name      = "CpuPercentage"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = var.cpuThreshold
  }

  criteria {
    metric_namespace = "Microsoft.Web/serverfarms"
    metric_name      = "DiskQueueLength"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = var.diskQueueLength
  }

  criteria {
    metric_namespace = "Microsoft.Web/serverfarms"
    metric_name      = "MemoryPercentage"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = var.memoryPercentage
  }

  criteria {
    metric_namespace = "Microsoft.Web/serverfarms"
    metric_name      = "HttpQueueLength"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = var.httpQueueLength
  }

  action {
    action_group_id = var.actionGroupId
  }
}

resource "azurerm_monitor_autoscale_setting" "autoscale" {
  count = var.kind == "FunctionApp" || var.kind == "elastic" ? 0 : var.instanceCount

  name                = "${local.baseName}-${count.index}-${var.location}"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = var.location
  target_resource_id  = azurerm_app_service_plan.sharedplan[count.index].id

  profile {
    name = "defaultProfile"

    capacity {
      default = var.environment == "prod" ? var.prodAutoScaleDefaultCapacity : 1
      minimum = var.environment == "prod" ? var.prodAutoScaleMinimumCapacity : 1
      maximum = var.environment == "prod" ? var.prodAutoScaleMaximumCapacity : 10
    }

    dynamic "rule" {
      for_each = length(var.autoscaleRules) > 0 ? var.autoscaleRules : local.defaultAutoscaleRules

      content {
        metric_trigger {
          metric_name        = rule.value["metricName"]
          metric_resource_id = azurerm_app_service_plan.sharedplan[count.index].id
          time_grain         = rule.value["timeGrain"]
          statistic          = rule.value["statistic"]
          time_window        = rule.value["timeWindow"]
          time_aggregation   = rule.value["timeAggregation"]
          operator           = rule.value["operator"]
          threshold          = rule.value["threshold"]
        }

        scale_action {
          direction = rule.value["direction"]
          type      = rule.value["type"]
          value     = rule.value["value"]
          cooldown  = rule.value["cooldown"]
        }
      }
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = var.autoScaleNotifySubscriptionAdmins
      send_to_subscription_co_administrator = var.autoScaleNotifyCoAdmins
      custom_emails                         = var.autoScaleNotifyEmails
    }
  }
}