variable "aws_region" {
  description = "AWS region to deploy into"
  type        = string
  default     = "us-east-1"
}

variable "client_name" {
  description = "Short client identifier used for resource naming, e.g. 'vetclinic-oak'"
  type        = string
}

variable "instance_type" {
  description = "EC2 instance type — t3.large recommended for full ShopStack"
  type        = string
  default     = "t3.large"
}

variable "disk_size_gb" {
  description = "Root volume size in GB"
  type        = number
  default     = 40
}

variable "key_pair_name" {
  description = "Name of existing EC2 key pair for SSH access"
  type        = string
}

variable "subnet_id" {
  description = "Subnet ID to launch the instance into"
  type        = string
}

variable "admin_cidr_blocks" {
  description = "CIDR blocks allowed SSH access (your management IP)"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}
