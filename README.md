Neurolog
================
> A Go-lang interface to access [`neural-redis`](https://github.com/antirez/neural-redis)

Documetation on [GoDoc](https://godoc.org/github.com/ziyasal/neurolog/neurolog)

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
docker run -it --link neural-redis:neurolog --rm reurolog redis-cli -h neurolog -p 6379
```

**Usage (Preview)**

```go
package main

import (
        "fmt"
        "time"
	"github.com/ziyasal/neurolog/neurolog"
)


func main() {

	options := neurolog.Options{
		Name:            "additions",
		Type:            "regressor",
		Inputs:          []string{"number1", "number2"},
		Outputs:         []string{"result"},
		HiddenLayers:    []int{3},
		DatasetSize:     50,
		TestDatasetSize: 10,
		RedisHost:       "localhost:6379",
	}

	network := neurolog.New(options)

    network.ObserveTrain(map[string]int64{"number1":3, "number2":5}, map[string]int64{"result":8})
    
    network.Train(0,0,true,true)
    
    for network.IsTraining() {
        fmt.Println("Training")
        time.Sleep(1)
    }
    
    fmt.Println(network.Run(map[string]int64{"number1":1, "number2":2}))
}


```
