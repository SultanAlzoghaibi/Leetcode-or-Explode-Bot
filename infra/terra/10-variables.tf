variable "db_password" {
  type        = string
  description = "The password for the database user"
  sensitive   = true
}
variable "db_user" {
  type        = string
  description = "Database username"
}
