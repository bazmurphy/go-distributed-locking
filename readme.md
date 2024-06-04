# Distributed Locking

run with: `docker compose up`

spawns one `etcd` and 3 `app` instances that all compete for a lock on a key to update it's value