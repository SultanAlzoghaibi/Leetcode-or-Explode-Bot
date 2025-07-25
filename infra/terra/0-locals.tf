locals {
  project_id = "lcc-backend"
  region = "us-central-1"
  apis = [
    "compute.googleapis.com",
    "container.googleapis.com",
    "logging.googleapis.com",
    "secretmanager.googleapis.com",
  ]
}