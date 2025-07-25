terraform {
  required_providers {
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "6.13.0"
    }
  }
}
resource "google_compute_subnetwork" "public" {
  name = "public"
  ip_cidr_range = "10.0.0.0/19" # search subnet calculator to see the range
  region = local.region
  network = google_compute_network.vpc.id
  private_ip_google_access = true
  stack_type = "IPV4_ONLY"
}

resource "google_compute_subnetwork" "private" {
  name    = "private"
  ip_cidr_range = "10.0.32.0/19"
  region = local.region
  network = google_compute_network.vpc.id
  private_ip_google_access = true
  stack_type = "IPV4_ONLY"

  secondary_ip_range {
    range_name = "k8s-pods"
    ip_cidr_range = "172.16.0.0/20" # 4096 pod ips availible
    # each time we take a port for a pod in k8s we use this range
  }
  secondary_ip_range {
    range_name = "k8s-services"
    ip_cidr_range = "172.20.0.0/24" # 256
  }
}