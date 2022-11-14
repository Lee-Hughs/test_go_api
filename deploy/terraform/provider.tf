provider "kubernetes" {
    config_path     = "~/.kube/config"
    config_context  = var.kube_context
}