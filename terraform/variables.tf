variable "project_id" {
  description = "The GCP Project ID"
  type        = string
}

variable "cloud_run_service_name" {
  description = "The name of the Cloud Run service"
  type        = string
  default     = "youtube-replacer-backend"
}

variable "region" {
  description = "The GCP region"
  type        = string
  default     = "us-east1"
}

variable "alert_email" {
  description = "Email address to send monitoring alerts to"
  type        = string
}
