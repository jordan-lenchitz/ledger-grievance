resource "kubernetes_config_map_v1" "mariadb_init" {
  metadata {
    name      = "mariadb-init-sql"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "mariadb"
    })
  }

  data = {
    "init.sql" = file("${path.module}/../db/init.sql")
  }
}

resource "kubernetes_service_v1" "mariadb" {
  metadata {
    name      = "mariadb-service"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "mariadb"
    })
  }

  spec {
    selector = {
      app = "mariadb"
    }

    port {
      name        = "mysql"
      protocol    = "TCP"
      port        = 3306
      target_port = 3306
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_stateful_set_v1" "mariadb" {
  metadata {
    name      = "mariadb"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels = merge(local.common_labels, {
      app = "mariadb"
    })
  }

  spec {
    service_name = kubernetes_service_v1.mariadb.metadata[0].name
    replicas     = 1

    selector {
      match_labels = {
        app = "mariadb"
      }
    }

    template {
      metadata {
        labels = merge(local.common_labels, {
          app = "mariadb"
        })
      }

      spec {
        container {
          name  = "mariadb"
          image = var.mariadb_image

          port {
            name           = "mysql"
            container_port = 3306
          }

          env_from {
            secret_ref {
              name = kubernetes_secret_v1.ledger_secrets.metadata[0].name
            }
          }

          resources {
            requests = {
              cpu    = "200m"
              memory = "256Mi"
            }
            limits = {
              cpu    = "1000m"
              memory = "1Gi"
            }
          }

          volume_mount {
            name       = "mariadb-data"
            mount_path = "/var/lib/mysql"
          }

          volume_mount {
            name       = "mariadb-init"
            mount_path = "/docker-entrypoint-initdb.d"
            read_only  = true
          }

          liveness_probe {
            exec {
              command = ["healthcheck.sh", "--connect", "--innodb_initialized"]
            }
            initial_delay_seconds = 15
            period_seconds        = 10
            timeout_seconds       = 5
          }

          readiness_probe {
            exec {
              command = ["healthcheck.sh", "--connect", "--innodb_initialized"]
            }
            initial_delay_seconds = 5
            period_seconds        = 5
            timeout_seconds       = 5
          }
        }

        volume {
          name = "mariadb-init"
          config_map {
            name = kubernetes_config_map_v1.mariadb_init.metadata[0].name
          }
        }
      }
    }

    volume_claim_template {
      metadata {
        name = "mariadb-data"
      }
      spec {
        access_modes = ["ReadWriteOnce"]
        resources {
          requests = {
            storage = var.db_storage_size
          }
        }
      }
    }
  }
}
