output "namespace" {
  description = "Kubernetes namespace hosting the ledger-grievance stack"
  value       = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
}

output "api_service_name" {
  description = "Kubernetes Service name for the ledger-grievance API"
  value       = kubernetes_service_v1.ledger_service.metadata[0].name
}

output "mariadb_service_name" {
  description = "Kubernetes Service name for MariaDB"
  value       = kubernetes_service_v1.mariadb.metadata[0].name
}

output "prometheus_service_name" {
  description = "Kubernetes Service name for Prometheus"
  value       = kubernetes_service_v1.prometheus.metadata[0].name
}

output "grafana_service_name" {
  description = "Kubernetes Service name for Grafana"
  value       = kubernetes_service_v1.grafana.metadata[0].name
}

output "grafana_admin_password" {
  description = "Generated or configured Grafana administrator password"
  value       = local.grafana_admin_password
  sensitive   = true
}
