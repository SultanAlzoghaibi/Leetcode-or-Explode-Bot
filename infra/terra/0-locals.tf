locals {
  project_id = "lcc-backend"
  region = "us-central1"
  zone =  "us-central1-a"
  apis = [
    "compute.googleapis.com",
    "container.googleapis.com",
    "logging.googleapis.com",
    "secretmanager.googleapis.com",
  ]
}