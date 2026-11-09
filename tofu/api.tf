resource "kubernetes_deployment_v1" "ledger_api" {
  metadata {
    name      = "ledger-api"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "ledger-api"
    })
  }

  spec {
    replicas = var.api_replicas

    selector {
      match_labels = {
        app = "ledger-api"
      }
    }

    template {
      metadata {
        labels = merge(local.common_labels, {
          app = "ledger-api"
        })
      }

      spec {
        container {
          name  = "ledger-api"
          image = var.api_image

          port {
            name           = "http"
            container_port = var.api_port
          }

          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.ledger_config.metadata[0].name
            }
          }

          env_from {
            secret_ref {
              name = kubernetes_secret_v1.ledger_secrets.metadata[0].name
            }
          }

          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
          }

          liveness_probe {
            http_get {
              path = "/health"
              port = var.api_port
            }
            initial_delay_seconds = 15
            period_seconds        = 20
          }

          readiness_probe {
            http_get {
              path = "/health/deep"
              port = var.api_port
            }
            initial_delay_seconds = 5
            period_seconds        = 10
          }
        }
      }
    }
  }

  depends_on = [
    kubernetes_stateful_set_v1.mariadb
  ]
}

resource "kubernetes_service_v1" "ledger_service" {
  metadata {
    name      = "ledger-service"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "ledger-api"
    })
  }

  spec {
    selector = {
      app = "ledger-api"
    }

    port {
      name        = "http"
      protocol    = "TCP"
      port        = 80
      target_port = var.api_port
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_horizontal_pod_autoscaler_v2" "ledger_hpa" {
  metadata {
    name      = "ledger-hpa"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "ledger-api"
    })
  }

  spec {
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment_v1.ledger_api.metadata[0].name
    }

    min_replicas = var.hpa_min_replicas
    max_replicas = var.hpa_max_replicas

    metric {
      type = "Resource"
      resource {
        name = "cpu"
        target {
          type                = "Utilization"
          average_utilization = var.hpa_cpu_utilization
        }
      }
    }
  }
}
