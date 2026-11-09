locals {
  common_labels = {
    "app.kubernetes.io/part-of"    = "ledger-grievance"
    "app.kubernetes.io/managed-by" = "opentofu"
    "environment"                  = var.environment
  }

  db_password            = var.db_password != "" ? var.db_password : random_password.db_password.result
  db_root_password       = var.db_root_password != "" ? var.db_root_password : random_password.db_root_password.result
  grafana_admin_password = var.grafana_admin_password != "" ? var.grafana_admin_password : random_password.grafana_admin.result
}

resource "random_password" "db_password" {
  length  = 24
  special = false
}

resource "random_password" "db_root_password" {
  length  = 24
  special = false
}

resource "random_password" "grafana_admin" {
  length  = 20
  special = false
}

resource "kubernetes_namespace_v1" "grievance_ledger" {
  metadata {
    name   = var.namespace
    labels = local.common_labels
  }
}

resource "kubernetes_config_map_v1" "ledger_config" {
  metadata {
    name      = "ledger-config"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels    = local.common_labels
  }

  data = {
    DB_HOST   = kubernetes_service_v1.mariadb.metadata[0].name
    DB_PORT   = "3306"
    DB_NAME   = var.db_name
    PORT      = tostring(var.api_port)
    LOG_LEVEL = var.log_level
  }
}

resource "kubernetes_secret_v1" "ledger_secrets" {
  metadata {
    name      = "ledger-secrets"
    namespace = kubernetes_namespace_v1.grievance_ledger.metadata[0].name
    labels    = local.common_labels
  }

  type = "Opaque"

  data = {
    DB_USER               = var.db_user
    DB_PASSWORD           = local.db_password
    MARIADB_ROOT_PASSWORD = local.db_root_password
    MARIADB_DATABASE      = var.db_name
    MARIADB_USER          = var.db_user
    MARIADB_PASSWORD      = local.db_password
  }
}
