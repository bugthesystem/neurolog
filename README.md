Neurolog
================
> A Go-lang interface to access `neural-redis`

**Build Docker image**
It contains `redis-server` and `neural-redis` pre-configured

```sh
docker built -t network .
```

**To run container**
```sh
docker run -d --name neural-redis -p 6379:6379 network
```

**To connect using `redis-cli`**
```sh
docker run -it --link neural-redis:network --rm network redis-cli -h network -p 6379
```

**Usage (Preview)**
```go
import "github.com/ziyasal/network/network"
```

```go
options := network.Options{
	Name:         "titanic",
	Type:         "classifier",
	Inputs:       []string{},
	Outputs:      []string{},
	HiddenLayers: []int{5},
	DatasetSize:  1000,
	TestDatasetSize:  500,
	RedisHost:    "localhost:6379",
}

network := network.New(options)
info := network.Info()

fmt.Println(info)
```
