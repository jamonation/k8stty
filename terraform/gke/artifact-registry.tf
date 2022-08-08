resource "google_artifact_registry_repository" "k8stty" {
  provider = google-beta
  project = var.project_id
  location = var.region
  repository_id = var.name
  description = "private k8stty docker repository with iam"
  format = "DOCKER"
}

resource "google_artifact_registry_repository_iam_binding" "serviceaccount-compute-iam-binding" {
  provider = google-beta
  project = google_artifact_registry_repository.k8stty.project
  location = google_artifact_registry_repository.k8stty.location
  repository = google_artifact_registry_repository.k8stty.name
  role = "roles/viewer"
    members = ["serviceAccount:${var.compute_engine_service_account}"]
}

resource "google_artifact_registry_repository_iam_member" "serviceaccount-compute-iam-member" {
  provider = google-beta
  project = var.project_id
  location = google_artifact_registry_repository.k8stty.location
  repository = google_artifact_registry_repository.k8stty.name
  role   = "roles/artifactregistry.reader"
  member = "serviceAccount:${var.compute_engine_service_account}"
}
