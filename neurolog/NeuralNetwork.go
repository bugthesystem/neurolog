package neurolog

import (
	"fmt"
	"reflect"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	redisPoolMaxIdle   = 80
	redisPoolMaxActive = 12000 // max number of connections
)

//NeuralNetwork ... Neural-Redis neuralNetwork
type NeuralNetwork struct {
	_opts  Options
	_key   string
	_flags map[string]string
	_rooms map[string]bool
	_pool  *redis.Pool
}

// Options ...
type Options struct {
	RedisHost    string
	Password     string
	Key          string
	Name         string
	Type         string
	Inputs       []string
	Outputs      []string
	HiddenLayers []int
	DatasetSize  int
	TestsetSize  int
	Normalize    bool
	RedisClient  *redis.Pool
	Prefix       string
	AutoCreate   bool
}

// New ... Creates new neuralNetwork using options
func New(opts Options) NeuralNetwork {
	if opts.Prefix == "" {
		opts.Prefix = "nn:"
	}

	neuralNetwork := NeuralNetwork{_opts: opts}

	if opts.RedisClient != nil {
		neuralNetwork._pool = opts.RedisClient
	}

	initRedisConnPool(&neuralNetwork, opts)

	return neuralNetwork
}

//Info ... Returns nn info
func (neurolog NeuralNetwork) Info() interface{} {
	c := neurolog._pool.Get()
	defer c.Close()

	fmt.Println(neurolog._opts.Name)
	n, err := c.Do("nr.info", neurolog._opts.Name)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}
	//TODO: zip result
	fmt.Println(reflect.TypeOf(n))
	fmt.Println(n)
	return n
}

func initRedisConnPool(neuralNetwork *NeuralNetwork, opts Options) {
	if opts.RedisHost == "" {
		panic("Missing redis `host`")
	}

	neuralNetwork._pool = newPool(opts)
}

func newPool(opts Options) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   redisPoolMaxIdle,
		MaxActive: redisPoolMaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", opts.RedisHost)
			if err != nil {
				return nil, err
			}

			if opts.Password != "" {
				if _, err := c.Do("AUTH", opts.Password); err != nil {
					c.Close()
					return nil, err
				}
				return c, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

}
