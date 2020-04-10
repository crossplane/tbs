package kubernetes.admission

# CloudSQLInstances must not be larger than 20GB.
deny[msg] {
    input.request.kind.kind = "CloudSQLInstance"
    input.request.operation = "CREATE"
    size := input.request.object.spec.forProvider.settings.dataDiskSizeGb
    size >= 20
    msg = sprintf("database size of %d GB is larger than limit of 20 GB", [size])
}