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
	RedisHost       string
	Password        string
	Key             string
	Name            string
	Type            string
	Inputs          []string
	Outputs         []string
	HiddenLayers    []int
	DatasetSize     int
	TestDatasetSize int
	Normalize       bool
	RedisClient     *redis.Pool
	Prefix          string
	AutoCreate      bool
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

//Info ... Returns Redis internal info about the neural network
func (network NeuralNetwork) Info() interface{} {
	c := network._pool.Get()
	defer c.Close()

	fmt.Println(network._opts.Name)
	n, err := c.Do("nr.info", network._opts.Name)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}
	//TODO: zip result
	fmt.Println(reflect.TypeOf(n))
	fmt.Println(n)
	return n
}

//Run ... Run the network returning a dict result
func (network NeuralNetwork) Run(input map[string]int64) map[string]string {

	_validateInput(input)

	inputLen := len(network._opts.Inputs)
	args := make([]string, inputLen+1)
	args[0] = network._opts.Name

	for i := 0; i < inputLen; i++ {
		args[i+1] = network._opts.Inputs[i]
	}

	fmt.Println(args)
	return map[string]string{}
}

func (network NeuralNetwork) ObserveTrain(input map[string]int64, output map[string]int64) bool {
	return network._observe(input, output, "train")
}

func (network NeuralNetwork) ObserveTest(input map[string]int64, output map[string]int64) bool {
	return network._observe(input, output, "test")
}

func (network NeuralNetwork) _observe(input map[string]int64, output map[string]int64, mode string) bool {
	return true
}

//Train ... Train
func (network NeuralNetwork) Train(maxCycles int, maxTime int, autoStop bool, backtrack bool) {
	//TODO:
}

func (network NeuralNetwork) IsTraining() bool {
	return true
}

func _validateInput(input map[string]int64) {
	//TODO:
}

func _validateOutput(output map[string]string) {
	//TODO:
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
