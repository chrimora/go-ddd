env "default" {
  dev = "docker://postgres"
  src = [
    "file://internal/post/infrastructure/sql/schema.sql",
    "file://internal/outbox/infrastructure/sql/schema.sql",
    "file://internal/user/infrastructure/sql/schema.sql",
  ]
}

