resource "google_sql_database" "database" {
  name     = "my-database"
  instance = google_sql_database_instance.instance.name
}


resource "google_compute_global_address" "private_ip_address" {
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc.id // we use our existinng vpc
  project       = local.project_id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc.id  // we use our existinng vpc
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

# See versions at https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database_instance#database_version
resource "google_sql_database_instance" "instance" {
  name             = "my-database-instance"
  region           = local.region
  database_version = "MYSQL_8_0"
  settings {
    tier = "db-f1-micro"

    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.vpc.self_link
    }
  }
  //todo set tru fro prod
  deletion_protection  = false
  depends_on = [
    google_service_networking_connection.private_vpc_connection
  ]

}



resource "google_sql_user" "default" {
  instance = google_sql_database_instance.instance.name
  name     = var.db_user
  password = var.db_password
}

resource "google_project_service" "sqladmin" {
  service = "sqladmin.googleapis.com"
}

resource "google_project_service" "compute" {
  service                     = "compute.googleapis.com"
  disable_on_destroy          = true
  disable_dependent_services  = true
}

resource "google_project_service" "servicenetworking" {
  service = "servicenetworking.googleapis.com"
}



