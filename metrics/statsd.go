package metrics

import (
	"fmt"
	"time"

	hystrixMetricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	hystrix "github.com/afex/hystrix-go/plugins"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

var statsD *statsd.Client

type StatsDConfig struct {
	Enabled     bool
	Host        string
	Port        int
	FlushPeriod time.Duration

	AppName string
}

func InitiateStatsDMetrics(config *StatsDConfig) error {
	const retryIntervalSeconds = 2
	const maxRetries = 5

	if config.Enabled {
		address := statsd.Address(fmt.Sprintf("%s:%d", config.Host, config.Port))
		prefix := statsd.Prefix(config.AppName)
		flushPeriod := statsd.FlushPeriod(config.FlushPeriod)

		var err error
		retries := 0
		for {
			statsD, err = statsd.New(address, prefix, flushPeriod)
			if err != nil {
				if retries < maxRetries {
					retries++
					time.Sleep(retryIntervalSeconds * time.Second)
					continue
				}
				return fmt.Errorf("error initiating statsD %+v", err)
			}
			break
		}

		hystrixCC := hystrix.StatsdCollectorConfig{
			StatsdAddr: fmt.Sprintf("%s:%d", config.Host, config.Port),
			Prefix:     config.AppName + ".hystrix",
		}
		collector, err := hystrix.InitializeStatsdCollector(&hystrixCC)
		if err != nil {
			statsD.Close()
			return fmt.Errorf("error initiating hystrix collector on statsD %+v", err)
		}
		hystrixMetricCollector.Registry.Register(collector.NewStatsdCollector)
	}
	return nil
}

func StatsDClient() *statsd.Client {
	return statsD
}

func CloseStatsDClient() {
	if statsD != nil {
		statsD.Close()
	}
}

func RecordHTTPStat(status int, path, method string, duration time.Duration) {
	const metric = "http.request,status=%d,path=%s"
	const metricTiming = "http.request.timing,status=%d,path=%s"

	statsD := StatsDClient()
	if statsD == nil {
		return
	}

	stat := fmt.Sprintf(metric, status, method+"-"+path)
	statsD.Increment(stat)

	statTiming := fmt.Sprintf(metricTiming, status, method+"-"+path)
	statsD.Timing(statTiming, int(duration/time.Millisecond))
}
