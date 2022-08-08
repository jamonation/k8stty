module "gke" {
  source                     = "terraform-google-modules/kubernetes-engine/google"
  project_id                 = var.project_id
  name                       = var.name
  region                     = var.region
  zones                      = var.zones
  network                    = "vpc-${var.name}"
  subnetwork                 = "${var.region}-${var.name}"
  ip_range_pods              = "${var.region}-${var.name}-pods"
  ip_range_services          = "${var.region}-${var.name}-services"
  http_load_balancing        = true
  network_policy             = true
  horizontal_pod_autoscaling = true
  filestore_csi_driver       = false
  service_account            = var.compute_engine_service_account
  default_max_pods_per_node  = 32
  remove_default_node_pool   = true

  node_pools = [
    {
      name                      = "${var.name}-pool"
      machine_type              = "${var.machine_type}"
      image_type		= "UBUNTU_CONTAINERD"
      #node_locations            = var.zones
      min_count                 = var.min_nodes
      max_count                 = var.max_nodes
      local_ssd_count           = 0
      disk_size_gb              = 100
      disk_type                 = "pd-standard"
      enable_gcfs               = false
      auto_repair               = true
      auto_upgrade              = true
      service_account           = var.compute_engine_service_account
      preemptible               = false
      initial_node_count        = 3
      max_pods_per_node  = 32
    },
  ]

  node_pools_oauth_scopes = {
    all = []

    "${var.name}-pool" = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }

  node_pools_labels = {
    all = {}

    "${var.name}-pool" = {
      "${var.name}-pool" = true
    }
  }

  node_pools_metadata = {
    all = {}

    "${var.name}-pool" = {
      node-pool-metadata-custom-value = "${var.name}-pool"
    }
  }

  node_pools_taints = {
    all = []

    "${var.name}-pool" = [
      {
        key    = "${var.name}-pool"
        value  = true
        effect = "PREFER_NO_SCHEDULE"
      },
    ]
  }

  node_pools_tags = {
    all = []

    "${var.name}-pool" = [
      "${var.name}-pool",
    ]
  }

  depends_on = [module.vpc]
}

module "vpc" {
    source  = "terraform-google-modules/network/google"
    version = "~> 4.0"

    project_id   = var.project_id
    network_name = "vpc-${var.name}"
    routing_mode = "GLOBAL"

    subnets = [
        {
            subnet_name           = "${var.region}-${var.name}"
            subnet_ip             = "10.254.0.0/23"
            subnet_region         = "${var.region}"
        },
    ]

    secondary_ranges = {
        "${var.region}-${var.name}" = [
            {
                range_name    = "${var.region}-${var.name}-services"
                ip_cidr_range = "10.254.2.0/23"
            },
             {
                range_name    = "${var.region}-${var.name}-pods"
                ip_cidr_range = "10.254.4.0/23"
            },
        ]
    }

}

resource "google_compute_global_address" "external-ip-address" {
  project      = var.project_id
  name         = "${var.name}-external-ip-address"
  address_type = "EXTERNAL"
  ip_version   = "IPV4"
}
