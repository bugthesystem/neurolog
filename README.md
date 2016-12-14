Neurolog
================
A Go-lang interface to access `neural-redis`

**Build Docker image**
It contains `redis-server` and `neural-redis` pre-configured

```sh
docker built -t neurolog .
```

**To run container**
```sh
docker run -d --name neural-redis neurolog
```

**To connect using `redis-cli`**
```sh
docker run -it --link neural-redis:neurolog --rm neurolog redis-cli -h neurolog -p 6379
```

**Usage**
```go
//TODO:
```
