env "default" {
  dev = "docker://postgres"
  src = [
    "file://internal/outbox/infrastructure/sql/schema.sql",
    "file://internal/user/infrastructure/sql/schema.sql",
    "file://internal/order/infrastructure/sql/schema.sql",
  ]
}

