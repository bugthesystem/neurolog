# Access from CLI

```sh
docker run -it --link some-redis:neurolog --rm neurolog redis-cli -h neurolog -p 6379
```