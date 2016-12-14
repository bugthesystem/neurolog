Neurolog
================
> A Go-lang interface to access `neural-redis`

**Build Docker image**
It contains `redis-server` and `neural-redis` pre-configured

```sh
docker built -t neurolog .
```

**To run container**
```sh
docker run -d --name neural-redis -p 6379:6379 neurolog
```

**To connect using `redis-cli`**
```sh
docker run -it --link neural-redis:neurolog --rm neurolog redis-cli -h neurolog -p 6379
```

**Usage (Preview)**
```go
import "github.com/ziyasal/neurolog/neurolog"
```

```go
    options := neurolog.Options{
		Name:         "titanic",
		Type:         "classifier",
		Inputs:       []string{},
		Outputs:      []string{},
		HiddenLayers: []int{5},
		DatasetSize:  1000,
		TestsetSize:  500,
		RedisHost:    "localhost:6379",
	}

	nn := neurolog.New(options)
	info := nn.Info()

	fmt.Println(info)
```
