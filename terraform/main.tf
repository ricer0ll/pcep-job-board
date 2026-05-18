terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.32.0"
    }
  }
}

provider "google" {
  project = "pcep-496700"
}

resource "google_secret_manager_secret" "discord-token" {
  secret_id = "discord-token"

  replication {
    auto {
    }
  }

  deletion_protection = false
}