variable "resource-group-name" {
  type = string
}

variable "appName" {
  type = string
}

variable "environment" {
  type = string
}

variable "location" {
  type = string
}

variable "autoScaleNotifyEmails" {
  type = list(string)
}