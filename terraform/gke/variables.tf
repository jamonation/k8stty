variable "name" {
  description = "General project name"
}

variable "project_id" {
  description = "The project ID to host the cluster in"
}

variable "compute_engine_service_account" {
  description = "GKE service account"
}

variable "region" {
  description = "Region for GKE cluster and gcloud resources"
}

variable "zones" {
  description = "Zones within gcloud region for nodes"
  type = list
}

variable "machine_type" {
  description = "Type of GCE machine for worker nodes"
}

variable "min_nodes" {
  description = "Minimum number of worker nodes"
}

variable "max_nodes" {
  description = "Maximum number of worker nodes"
}

variable "storage_bucket_prefix" {
  description = "Path to gcs storage location"
}
