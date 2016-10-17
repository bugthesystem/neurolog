# Access from CLI

```sh
docker run -it --link some-redis:hasan --rm hasan redis-cli -h hasan -p 6379
```