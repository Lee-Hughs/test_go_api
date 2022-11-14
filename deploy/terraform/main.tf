resource "kubernetes_namespace" "example" {
    metadata {
        name = "tf-test-lee"
    }
}