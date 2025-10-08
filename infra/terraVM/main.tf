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

variable "username" {
  default = "sultanalzoghaibi"
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
    apt-get install -y git docker.io docker-compose
    usermod -aG docker $(whoami)

    # Copy project files from metadata archive (uploaded below)
    cd /home/${var.username}
    tar -xzf /tmp/project.tar.gz
    cd docker-compose
    docker compose up -d
  EOT

  # 👇 Upload your docker-compose folder as metadata for startup use
  metadata = {
    "gce-container-declaration" = ""
  }

  provisioner "file" {
    source      = "../Docker/docker-compose"
    destination = "/tmp/project"
  }

  provisioner "remote-exec" {
    inline = [
      "cd /tmp",
      "tar -czf project.tar.gz project",
      "mv project.tar.gz /tmp/",
      "rm -rf /tmp/project"
    ]

    connection {
      type        = "ssh"
      user        = "sultanalzoghaibi"
      private_key = file("~/.ssh/id_rsa")
      host        = google_compute_instance.lcc_vm_new.network_interface[0].access_config[0].nat_ip
    }
  }
}