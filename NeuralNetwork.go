package neurolog

import (
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
	Host         string
	Password     string
	Key          string
	Name         string
	Type         string
	Inputs       []string
	Outputs      []string
	HiddenLayers []string
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

	//TODO

	return neuralNetwork
}

func info(neurolog NeuralNetwork) map[string]string {
	result := map[string]string{}
	c := neurolog._pool.Get()
	defer c.Close()

	n, err := c.Do("nr.info", neurolog._opts.Name)
	if err != nil {
		//TODO: handle error return from c.Do or type conversion error.
	}
	//dict(zip(result[0::2], result[1::2]))
	return result
}

func initRedisConnPool(neuralNetwork *NeuralNetwork, opts Options) {
	if opts.Host == "" {
		panic("Missing redis `host`")
	}

	neuralNetwork._pool = newPool(opts)
}

func newPool(opts Options) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   redisPoolMaxIdle,
		MaxActive: redisPoolMaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", opts.Host)
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
