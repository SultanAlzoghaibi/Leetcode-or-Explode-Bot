resource "google_container_cluster" "gke" {
  name = "lcc-backend"
  location = local.region
  remove_default_node_pool = true
  initial_node_count = 1
  network = google_compute_network.vpc.self_link
  subnetwork = google_compute_subnetwork.private.self_link
  networking_mode = "VPC_NATIVE"

  deletion_protection = false //todo change this in prod
  addons_config {
    http_load_balancing {
      disabled = true
    }
    horizontal_pod_autoscaling {
      disabled = true
    }

  }
  release_channel {
    channel = "REGULAR"
  }

  workload_identity_config {
    workload_pool = "${local.project_id}.svc.id.goog"
  }

  ip_allocation_policy {
    cluster_secondary_range_name = "k8s-pods"
    services_secondary_range_name = "k8s-services"
  }
  private_cluster_config {
    enable_private_nodes = true
    enable_private_endpoint = false
    master_ipv4_cidr_block = "192.168.0.0/28"
  }

}
resource "google_compute_address" "static_ip" {
  name         = "static-ip"
  address_type = "EXTERNAL"
  region       = local.region
}
