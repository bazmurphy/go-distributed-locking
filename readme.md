Pull the etcd image

```sh
docker pull bitnami/etcd:latest
```

Run a container with it

```sh
docker run -it --rm -p 2379:2379 --env ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd:latest
```