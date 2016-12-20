package neurolog

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

const (
	redisPoolMaxIdle   = 80
	redisPoolMaxActive = 12000 // max number of connections
)

//NeuralNetwork ... Neural-Redis interface
type NeuralNetwork struct {
	_opts Options
	_pool *redis.Pool
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
func (network NeuralNetwork) Info() map[string]string {
	c := network._pool.Get()
	defer c.Close()

	cmdResult, err := c.Do("nr.info", network._opts.Name)
	if err != nil {
		log.Error(err)
	}

	result := map[string]string{}

	if resultArray, ok := cmdResult.([]interface{}); ok {
		for i := 0; i < len(resultArray); i += 2 {
			result[toString(resultArray[i])] = toString(resultArray[i+1])
		}
	}

	return result
}

//Classify ...  Run the network returning the classified class
func (network NeuralNetwork) Classify(input map[string]int64) interface{} {

	_validateInput(network, input)

	args := []interface{}{network._opts.Name}

	for _, key := range network._opts.Inputs {
		args = append(args, input[key])
	}

	fmt.Println(args)

	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("nr.class", args...)
	if err != nil {
		log.Error(err)
	}

	if str, ok := result.(string); ok {
		idx, _ := strconv.ParseInt(string(str), 10, 64)
		return network._opts.Outputs[idx]
	}

	return map[string]int64{}

}

//Run ... Run the network returning a dict result
func (network NeuralNetwork) Run(input map[string]int64) map[string]float64 {

	_validateInput(network, input)

	args := []interface{}{network._opts.Name}

	for _, key := range network._opts.Inputs {
		args = append(args, input[key])
	}

	fmt.Println(args)

	c := network._pool.Get()
	defer c.Close()

	cmdResult, err := c.Do("nr.run", args...)

	if err != nil {
		log.Error(err)
	}

	result := map[string]float64{}

	if resultArray, ok := cmdResult.([]interface{}); ok {
		for i := 0; i < len(network._opts.Outputs); i++ {
			if values, ok := resultArray[i].([]uint8); ok {
				floatValue, _ := strconv.ParseFloat(string(values[:]), 64)
				result[network._opts.Outputs[i]] = floatValue
			} else {
				result[network._opts.Outputs[i]] = 0
			}

		}
	}

	return result
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

	args := []interface{}{network._opts.Name}

	for _, key := range network._opts.Inputs {
		args = append(args, input[key])
	}

	args = append(args, "->")

	if strings.ToLower(network._opts.Type) == "regressor" {
		for _, key := range network._opts.Outputs {
			args = append(args, output[key])
		}
	} else {
		clsType := 0
		for i, name := range network._opts.Outputs {

			if _, ok := output[name]; ok {
				clsType = i
				break
			}

		}
		args = append(args, clsType)

	}

	args = append(args, mode)

	fmt.Println(args)

	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("nr.observe", args...)
	if err != nil {
		log.Error(err)
	}

	return result
}

//IsCreated ... Returns true if the neural network is created
func (network NeuralNetwork) IsCreated() bool {
	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("EXISTS", network._opts.Name)
	if err != nil {
		log.Error(err)
	}

	if parseResult, ok := result.(int64); ok {
		return parseResult == 1
	} else {
		return false
	}

}

//Create ... Create the neural network
func (network NeuralNetwork) Create() int64 {

	args := []interface{}{network._opts.Name, network._opts.Type, strconv.Itoa(len(network._opts.Inputs))}

	for _, el := range network._opts.HiddenLayers {
		args = append(args, el)
	}

	args = append(args, "->")
	args = append(args, len(network._opts.Outputs))

	if network._opts.Normalize {
		args = append(args, "NORMALIZE")
	}
	args = append(args, "DATASET")

	args = append(args, network._opts.DatasetSize)
	args = append(args, "TEST")
	args = append(args, network._opts.TestDatasetSize)

	fmt.Println(args)

	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("nr.create", args...)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(result)

	if flag, ok := result.(int64); ok {
		return flag
	} else {
		return 0
	}
}

//Delete ...  Delete the neural network
func (network NeuralNetwork) Delete() {
	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("DEL", network._opts.Name)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(result)
}

//Re-create ... Recreate the neural network
func (network NeuralNetwork) ReCreate() {
	network.Delete()
	network.Create()
}

//Train ... Train
func (network NeuralNetwork) Train(maxCycles int, maxTime int, autoStop bool, backtrack bool) interface{} {
	args := []interface{}{network._opts.Name}

	if maxCycles != 0 {
		args = append(args, "MAXCYCLES")
		args = append(args, maxCycles)
	}

	if maxTime != 0 {
		args = append(args, "MAXTIME")
		args = append(args, maxTime)
	}

	if autoStop {
		args = append(args, "AUTOSTOP")
	}

	if backtrack {
		args = append(args, "BACKTRACK")
	}

	c := network._pool.Get()
	defer c.Close()

	result, err := c.Do("nr.train", args...)
	if err != nil {
		log.Error(err)
	}

	return result
}

func (network NeuralNetwork) IsTraining() bool {
	info := network.Info()
	return info["training"] == "1"
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
