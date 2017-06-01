package neurolog

import "github.com/garyburd/redigo/redis"

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
