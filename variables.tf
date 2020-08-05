variable "resource-group-name" {
  type        = string
  description = "Name of the resource group that exists in Azure"
}

variable "appName" {
  type        = string
  description = "The base name of the application used in the naming convention."
}

variable "environment" {
  type        = string
  description = "Name of the environment ex (Dev, Test, QA, Prod)"
}

variable "location" {
  type        = string
  description = "(Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created."
}

variable "instanceCount" {
  type        = number
  description = "The number of instances you want to create in this location"
  default     = 1
}

variable "kind" {
  type        = string
  description = "The kind of the App Service Plan to create. Possible values are Windows (also available as App), Linux, elastic (for Premium Consumption) and FunctionApp (for a Consumption Plan). Defaults to Windows. Changing this forces a new resource to be created."
  default     = "Windows"
}

variable "reserved" {
  type        = bool
  description = "Is this App Service Plan Reserved. Defaults to false. NOTE: When creating a Linux App Service Plan, the reserved field must be set to true, and when creating a Windows/app App Service Plan the reserved field must be set to false."
  default     = false
}

variable "skuTier" {
  type        = string
  description = "(Required) Specifies the plan's pricing tier. https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-manager-sku-not-available-errors"
  default     = "Standard"
}

variable "skuSize" {
  type        = string
  description = "(Required) Specifies the plan's instance size. https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-manager-sku-not-available-errors"
  default     = "S1"
}

variable "skuCapacity" {
  type        = number
  description = "(Optional) Specifies the number of workers associated with this App Service Plan."
  default     = 1
}

variable "cpuThreshold" {
  type        = number
  description = "CPU metric alerts for an app service plan"
  default     = 70
}

variable "diskQueueLength" {
  type        = number
  description = "Disk Queue Length metric alerts for an app service plan"
  default     = 100
}

variable "memoryPercentage" {
  type        = number
  description = "Memory metric alerts for an app service plan"
  default     = 90
}

variable "httpQueueLength" {
  type        = number
  description = "Http Queue Length metric alerts for an app service plan"
  default     = 100
}

variable "actionGroupId" {
  type        = string
  description = "The ID for the action group to receive notifications of alerts"
}

variable "prodAutoScaleDefaultCapacity" {
  type        = number
  description = "An optional override of the production default capacity for auto scaling for non consumption plans"
  default     = 2
}

variable "prodAutoScaleMinimumCapacity" {
  type        = number
  description = "An optional override of the production minimum capacity for auto scaling for non consumption plans"
  default     = 2
}

variable "prodAutoScaleMaximumCapacity" {
  type        = number
  description = "An optional override of the production maximum capacity for auto scaling for non consumption plans"
  default     = 10
}

variable "autoScaleNotifySubscriptionAdmins" {
  type        = bool
  description = "Turn on or off email notifications to subscription admins when auto scaling triggers"
  default     = true
}

variable "autoScaleNotifyCoAdmins" {
  type        = bool
  description = "Turn on or off email notifications to co admins when auto scaling triggers"
  default     = true
}

variable "autoScaleNotifyEmails" {
  type        = list(string)
  description = "A list of emails to notify when auto scaling triggers"
}

variable "autoscaleRules" {
  type        = list(map(string))
  description = "An optional override of the auto scale rules. Use Default for example structure to override."
  default     = []
}