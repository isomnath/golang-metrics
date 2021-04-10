package metrics

import (
	raven "github.com/getsentry/raven-go"
)

var client *raven.Client

type SentryConfig struct {
	Enabled bool
	DSN     string
}

func InitSentry(config *SentryConfig) error {
	if config.Enabled {
		var err error
		client, err = raven.New(config.DSN)
		if err != nil {
			return err
		}
	}
	return nil
}

func StopSentry() {
	if client != nil {
		client.Close()
	}
}

func CaptureError(err error) {
	if client != nil {
		client.CaptureError(err, map[string]string{})
	}
}

func CaptureWarn(err error) {
	if client != nil {
		client.CaptureError(err, warnTags)
	}
}

func CaptureErrorWithTags(err error, tags map[string]string) {
	if client != nil {
		client.CaptureError(err, tags)
	}
}

var warnTags = map[string]string{
	"level": string(raven.WARNING),
}
