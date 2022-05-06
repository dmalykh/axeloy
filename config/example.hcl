
database "pgx" {
  dsn = "postgre://name@pass:serv"
}

driver "graphql" "ways/graphql/graphql.so" {
  listen_addr = "127.0.0.1:8080"
}

driver "superdemo" "ways/graphql/graphql.so" {
  port = "998"
}