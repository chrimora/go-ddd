package config

import "time"

type ForwarderConfig struct {
	MaxRetries         int           `env:"MAX_RETRIES" envDefault:"5"`
	PublisherBatchSize int           `env:"PUBLISHER_BATCH" envDefault:"10"`
	PublisherSleep     time.Duration `env:"PUBLISHER_SLEEP" envDefault:"300ms"`
	WatchdogTick       time.Duration `env:"WATCHDOG_TICK" envDefault:"3s"`
	WatchdogStaleLimit time.Duration `env:"WATCHDOG_STALE_LIMIT" envDefault:"3s"`
}

type RouterConfig struct {
	MaxRetries              int           `env:"MAX_RETRIES" envDefault:"5"`
	RetryInterval           time.Duration `env:"RETRY_INTERVAL" envDefault:"1s"`
	RetryIntervalMultiplier float64       `env:"RETRY_INTERVAL_MULT" envDefault:"2"`
}
