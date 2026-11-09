variable "namespace" {
  description = "Kubernetes namespace for the ledger-grievance stack"
  type        = string
  default     = "grievance-ledger"
}

variable "environment" {
  description = "Deployment environment label (e.g., production, staging, local)"
  type        = string
  default     = "production"
}

variable "kubeconfig_path" {
  description = "Path to the kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "kube_context" {
  description = "Optional kubeconfig context name to target"
  type        = string
  default     = ""
}

variable "api_image" {
  description = "Container image for the ledger-grievance Go API"
  type        = string
  default     = "jordan-lenchitz/ledger-grievance-api:latest"
}

variable "api_replicas" {
  description = "Initial number of replicas for the ledger-api deployment"
  type        = number
  default     = 2
}

variable "api_port" {
  description = "Port exposed by the ledger-grievance HTTP server"
  type        = number
  default     = 8000
}

variable "log_level" {
  description = "Structured log level for the API (debug, info, warn, error)"
  type        = string
  default     = "info"
}

variable "hpa_min_replicas" {
  description = "Minimum number of replicas for HorizontalPodAutoscaler"
  type        = number
  default     = 2
}

variable "hpa_max_replicas" {
  description = "Maximum number of replicas for HorizontalPodAutoscaler"
  type        = number
  default     = 10
}

variable "hpa_cpu_utilization" {
  description = "Target average CPU utilization percentage for autoscaling kindness"
  type        = number
  default     = 70
}

variable "mariadb_image" {
  description = "Container image for MariaDB"
  type        = string
  default     = "mariadb:11.4"
}

variable "db_name" {
  description = "Database name for storing grievances"
  type        = string
  default     = "grievance_ledger"
}

variable "db_user" {
  description = "Database application username"
  type        = string
  default     = "ledger"
}

variable "db_password" {
  description = "Optional explicit database password; if empty, a secure random password is generated"
  type        = string
  default     = ""
  sensitive   = true
}

variable "db_root_password" {
  description = "Optional explicit MariaDB root password; if empty, a secure random password is generated"
  type        = string
  default     = ""
  sensitive   = true
}

variable "db_storage_size" {
  description = "Persistent volume storage request for MariaDB"
  type        = string
  default     = "10Gi"
}

variable "prometheus_image" {
  description = "Container image for Prometheus"
  type        = string
  default     = "prom/prometheus:v2.53.0"
}

variable "grafana_image" {
  description = "Container image for Grafana"
  type        = string
  default     = "grafana/grafana:11.0.0"
}

variable "grafana_url" {
  description = "URL for the Grafana provider when managing standalone/port-forwarded Grafana instances"
  type        = string
  default     = "http://localhost:3000"
}

variable "grafana_admin_user" {
  description = "Grafana administrator username"
  type        = string
  default     = "admin"
}

variable "grafana_admin_password" {
  description = "Optional explicit Grafana admin password; if empty, a secure random password is generated"
  type        = string
  default     = ""
  sensitive   = true
}

variable "enable_grafana_provider_resources" {
  description = "Whether to also provision datasources/dashboards directly via the Grafana HTTP API provider (useful against local docker-compose Grafana)"
  type        = bool
  default     = false
}
