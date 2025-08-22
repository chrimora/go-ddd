env "default" {
  dev = "docker://postgres"
  src = [
    "file://internal/domain/user/schema.sql",
    "file://internal/outbox/schema.sql",
  ]
}

