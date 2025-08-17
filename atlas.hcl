env "default" {
  dev = "docker://postgres"
  src = [
    "file://user/schema.sql",
  ]
}

