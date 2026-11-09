resource "kubernetes_config_map_v1" "prometheus_config" {
  metadata {
    name      = "prometheus-config"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "prometheus"
    })
  }

  data = {
    "prometheus.yml" = <<-EOT
      global:
        scrape_interval: 15s

      scrape_configs:
        - job_name: 'ledger-grievance'
          static_configs:
            - targets: ['${kubernetes_service_v1.ledger_service.metadata[0].name}:80']
    EOT
  }
}

resource "kubernetes_deployment_v1" "prometheus" {
  metadata {
    name      = "grievance-prometheus"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "prometheus"
    })
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "prometheus"
      }
    }

    template {
      metadata {
        labels = merge(local.common_labels, {
          app = "prometheus"
        })
      }

      spec {
        container {
          name  = "prometheus"
          image = var.prometheus_image

          port {
            name           = "http"
            container_port = 9090
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "256Mi"
            }
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }

          volume_mount {
            name       = "config"
            mount_path = "/etc/prometheus/prometheus.yml"
            sub_path   = "prometheus.yml"
            read_only  = true
          }
        }

        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map_v1.prometheus_config.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service_v1" "prometheus" {
  metadata {
    name      = "prometheus-service"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "prometheus"
    })
  }

  spec {
    selector = {
      app = "prometheus"
    }

    port {
      name        = "http"
      protocol    = "TCP"
      port        = 9090
      target_port = 9090
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_secret_v1" "grafana_secrets" {
  metadata {
    name      = "grafana-secrets"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "grafana"
    })
  }

  type = "Opaque"

  data = {
    GF_SECURITY_ADMIN_USER     = var.grafana_admin_user
    GF_SECURITY_ADMIN_PASSWORD = local.grafana_admin_password
  }
}

resource "kubernetes_config_map_v1" "grafana_datasources" {
  metadata {
    name      = "grafana-datasources"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "grafana"
    })
  }

  data = {
    "prometheus.yaml" = <<-EOT
      apiVersion: 1
      datasources:
        - name: Prometheus
          type: prometheus
          uid: prometheus-ledger
          access: proxy
          url: http://${kubernetes_service_v1.prometheus.metadata[0].name}:9090
          isDefault: true
          editable: false
    EOT
  }
}

resource "kubernetes_config_map_v1" "grafana_dashboard_provider" {
  metadata {
    name      = "grafana-dashboard-provider"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "grafana"
    })
  }

  data = {
    "dashboards.yaml" = <<-EOT
      apiVersion: 1
      providers:
        - name: 'wholesome-dashboards'
          orgId: 1
          folder: 'Ledger Grievance'
          type: file
          disableDeletion: false
          editable: true
          options:
            path: /var/lib/grafana/dashboards
    EOT
  }
}

resource "kubernetes_config_map_v1" "grafana_wholesome_dashboard" {
  metadata {
    name      = "grafana-wholesome-dashboard"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "grafana"
    })
  }

  data = {
    "wholesome.json" = file("${path.module}/dashboards/wholesome.json")
  }
}

resource "kubernetes_deployment_v1" "grafana" {
  metadata {
    name      = "grievance-grafana"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "grafana"
    })
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "grafana"
      }
    }

    template {
      metadata {
        labels = merge(local.common_labels, {
          app = "grafana"
        })
      }

      spec {
        container {
          name  = "grafana"
          image = var.grafana_image

          port {
            name           = "http"
            container_port = 3000
          }

          env_from {
            secret_ref {
              name = kubernetes_secret_v1.grafana_secrets.metadata[0].name
            }
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }

          volume_mount {
            name       = "datasources"
            mount_path = "/etc/grafana/provisioning/datasources"
            read_only  = true
          }

          volume_mount {
            name       = "dashboard-provider"
            mount_path = "/etc/grafana/provisioning/dashboards"
            read_only  = true
          }

          volume_mount {
            name       = "dashboards"
            mount_path = "/var/lib/grafana/dashboards"
            read_only  = true
          }
        }

        volume {
          name = "datasources"
          config_map {
            name = kubernetes_config_map_v1.grafana_datasources.metadata[0].name
          }
        }

        volume {
          name = "dashboard-provider"
          config_map {
            name = kubernetes_config_map_v1.grafana_dashboard_provider.metadata[0].name
          }
        }

        volume {
          name = "dashboards"
          config_map {
            name = kubernetes_config_map_v1.grafana_wholesome_dashboard.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service_v1" "grafana" {
  metadata {
    name      = "grafana-service"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "grafana"
    })
  }

  spec {
    selector = {
      app = "grafana"
    }

    port {
      name        = "http"
      protocol    = "TCP"
      port        = 3000
      target_port = 3000
    }

    type = "ClusterIP"
  }
}

resource "grafana_data_source" "prometheus" {
  count = var.enable_grafana_provider_resources ? 1 : 0

  type       = "prometheus"
  name       = "Prometheus (OpenTofu)"
  uid        = "prometheus-ledger"
  url        = "http://prometheus:9090"
  is_default = true
}

resource "grafana_dashboard" "wholesome" {
  count = var.enable_grafana_provider_resources ? 1 : 0

  config_json = file("${path.module}/dashboards/wholesome.json")
  overwrite   = true

  depends_on = [
    grafana_data_source.prometheus
  ]
}
