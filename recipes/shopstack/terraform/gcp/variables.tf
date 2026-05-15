variable "gcp_project" {
  description = "GCP project ID"
  type        = string
}

variable "gcp_region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "gcp_zone" {
  description = "GCP zone"
  type        = string
  default     = "us-central1-a"
}

variable "client_name" {
  description = "Short client identifier used for resource naming, e.g. 'vetclinic-oak'"
  type        = string
}

variable "machine_type" {
  description = "GCP machine type — e2-standard-2 recommended for full ShopStack"
  type        = string
  default     = "e2-standard-2"
}

variable "disk_size_gb" {
  description = "Boot disk size in GB"
  type        = number
  default     = 40
}

variable "ssh_public_key" {
  description = "SSH public key content for the admin user"
  type        = string
}

variable "admin_cidr_blocks" {
  description = "CIDR blocks allowed SSH access (your management IP)"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}
