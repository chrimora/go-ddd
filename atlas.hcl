env "default" {
  dev = "docker://postgres"
  src = [
    "file://internal/infrastructure/sql/user/schema.sql",
    "file://internal/infrastructure/sql/outbox/schema.sql",
  ]
}

