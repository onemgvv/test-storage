package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

const (
	defaultHttlPort               = "5535"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	defaultLimiterRPS   = 10
	defaultLimiterBurst = 20
	defaultLimiterTTL   = 10 * time.Minute
)

type (
	Config struct {
		HTTP    HTTPConfig    `json:"http"`
		Limiter LimiterConfig `json:"limiter"`
		Storage StorageConfig `json:"storage"`
	}

	HTTPConfig struct {
		Port               string         `json:"port"`
		Timeouts           TimeoutsConfig `json:"timeouts"`
		MaxHeaderMegabytes int            `json:"maxHeaderMegabytes"`
	}

	TimeoutsConfig struct {
		Read  time.Duration `json:"read"`
		Write time.Duration `json:"write"`
	}

	LimiterConfig struct {
		RPS   int           `json:"rps"`
		Burst int           `json:"burst"`
		TTL   time.Duration `json:"ttl"`
	}

	StorageConfig struct {
		DefaultTTL time.Duration `json:"defaultTTL"`
	}
)

func (c *Config) Init(path string) error {
	*c = setDefaultValues()

	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(jsonFile, &c); err != nil {
		return err
	}

	c.Limiter.TTL *= time.Minute
	c.HTTP.Timeouts.Read *= time.Second
	c.HTTP.Timeouts.Write *= time.Second
	c.Storage.DefaultTTL *= time.Minute

	return nil
}

func setDefaultValues() Config {
	timeoutsConf := TimeoutsConfig{
		Read:  defaultHTTPRWTimeout,
		Write: defaultHTTPRWTimeout,
	}

	httpConf := HTTPConfig{
		Port:               defaultHttlPort,
		Timeouts:           timeoutsConf,
		MaxHeaderMegabytes: defaultHTTPMaxHeaderMegabytes,
	}

	limiterConf := LimiterConfig{
		RPS:   defaultLimiterRPS,
		Burst: defaultLimiterBurst,
		TTL:   defaultLimiterTTL,
	}

	return Config{
		HTTP:    httpConf,
		Limiter: limiterConf,
	}
}
