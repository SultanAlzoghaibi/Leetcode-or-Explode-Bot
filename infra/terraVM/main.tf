terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
}

provider "google" {
  project = "lcc-backend"
  region  = "us-central1"
  zone    = "us-central1-a"
}

resource "google_compute_instance" "lcc_vm_new" {
  name         = "lcc-vm-new"
  machine_type = "e2-micro"
  zone         = "us-central1-a"

  tags = ["allow-ssh", "https-server"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
      size  = 10
    }
  }

  network_interface {
    network    = "projects/lcc-backend/global/networks/main"
    subnetwork = "projects/lcc-backend/regions/us-central1/subnetworks/public"

    # Attach your reserved external IP
    access_config {
      nat_ip = "34.170.46.132"
    }
  }

  metadata_startup_script = <<-EOT
    #!/bin/bash
    apt-get update -y
    apt-get install -y docker.io docker-compose
    usermod -aG docker $(whoami)
    systemctl enable docker
    systemctl start docker
  EOT
}