package neurolog

import (
	"fmt"
	"reflect"
	"strconv"
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

	if !neuralNetwork.IsCreated() {
		neuralNetwork.Create()
	}

	return neuralNetwork
}

//Info ... Returns Redis internal info about the neural network
func (network NeuralNetwork) Info() interface{} {
	c := network._pool.Get()
	defer c.Close()

	fmt.Println(network._opts.Name)
	result, err := c.Do("nr.info", network._opts.Name)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}

	//TODO: zip result
	fmt.Println(reflect.TypeOf(result))
	fmt.Println(result)

	//dict(zip(result[0::2], result[1::2]))
	return map[string]string{}
}

//Classify ...  Run the network returning the classified class
func (network NeuralNetwork) Classify(input map[string]int64) interface{} {

	_validateInput(network, input)

	args := []string{network._opts.Name}

	for i := 0; i < len(network._opts.Inputs); i++ {
		args = append(args, network._opts.Inputs[i])
	}

	fmt.Println(args)

	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("nr.class", args)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}

	if str, ok := result.(string); ok {
		idx, _ := strconv.ParseInt(string(str), 10, 64)
		return network._opts.Outputs[idx]
	} else {
		return nil
	}

}

//Run ... Run the network returning a dict result
func (network NeuralNetwork) Run(input map[string]int64) map[string]string {

	_validateInput(network, input)

	args := []string{network._opts.Name}

	for i := 0; i < len(network._opts.Inputs); i++ {
		args = append(args, network._opts.Inputs[i])
	}

	fmt.Println(args)
	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("nr.run", args)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}

	fmt.Println(result)

	//dict(zip(self.outputs, result))
	return map[string]string{}
}

//ObserveTrain ...  Add a data sample into the training dataset
func (network NeuralNetwork) ObserveTrain(input map[string]int64, output map[string]int64) interface{} {
	return network._observe(input, output, "train")
}

//ObserveTest ...  Add a data sample into the testing dataset
func (network NeuralNetwork) ObserveTest(input map[string]int64, output map[string]int64) interface{} {
	return network._observe(input, output, "test")
}

//Add a data sample into the training or testing dataset
//Mode can be `train` or `test`
func (network NeuralNetwork) _observe(input map[string]int64, output map[string]int64, mode string) interface{} {

	_validateInput(network, input)
	_validateOutput(network, output)

	args := []string{network._opts.Name}

	for i := 0; i < len(network._opts.Inputs); i++ {
		args = append(args, network._opts.Inputs[i])
	}

	args = append(args, "->")

	if network._opts.Type == "regressor" {
		for i := 0; i < len(network._opts.Outputs); i++ {
			args = append(args, network._opts.Inputs[i])
		}
	} else {
		clsType := 0
		for i, name := range network._opts.Outputs {

			if _, ok := output[name]; ok {
				clsType = i
				break
			}

		}
		args = append(args, string(clsType))

	}
	args = append(args, mode)

	fmt.Println(args)

	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("nr.observe", args)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}

	return result
}

//IsCreated ... Returns true if the neural network is created
func (network NeuralNetwork) IsCreated() bool {
	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("EXISTS", network._opts.Name)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}

	if str, ok := result.(string); ok {
		return string(str) == "1"
	} else {
		return false
	}

}

//Create ... Create the neural network
func (network NeuralNetwork) Create() interface{} {

	args := []string{network._opts.Name, network._opts.Type, string(len(network._opts.Inputs))}
	for _, el := range network._opts.HiddenLayers {
		args = append(args, string(el))
	}

	args = append(args, "->")
	args = append(args, string(len(network._opts.Outputs)))

	if network._opts.Normalize {
		args = append(args, "NORMALIZE")
	}
	args = append(args, "DATASET")

	args = append(args, string(network._opts.DatasetSize))
	args = append(args, "TEST")
	args = append(args, string(network._opts.TestDatasetSize))

	c := network._pool.Get()
	defer c.Close()
	result, err := c.Do("nr.create", args)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}
	fmt.Println(result)
	return result
}

//Delete ...  Delete the neural network
func (network NeuralNetwork) Delete() {
	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("DELETE", network._opts.Name)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}
	fmt.Println(result)
}

//ReCreate ... Recreate the neural network
func (network NeuralNetwork) ReCreate() {
	network.Delete()
	network.Create()
}

//Train ... Train
func (network NeuralNetwork) Train(maxCycles int, maxTime int, autoStop bool, backtrack bool) interface{} {
	args := []string{network._opts.Name}

	if maxCycles != 0 {
		args = append(args, "MAXCYCLES")
		args = append(args, string(maxCycles))
	}

	if maxTime != 0 {
		args = append(args, "MAXTIME")
		args = append(args, string(maxTime))
	}

	if autoStop {
		args = append(args, "AUTOSTOP")
	}

	if autoStop {
		args = append(args, "BACKTRACK")
	}

	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("nr.train", args)
	if err != nil {
		fmt.Println(err)
		//TODO: handle error return from c.Do or type conversion error.
	}
	fmt.Println(result)
	return result
}

func (network NeuralNetwork) IsTraining() bool {
	info := network.Info()
	return info["training"] == 1
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
