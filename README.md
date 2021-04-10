# Golang Metrics Collector

This Library Exposes APIs to publish metrics for the following monitoring platforms:

* **New Relic**
* **Sentry**
* **StatsD**

For each platform, the usage can be controlled via a feature toggle.

## Sample Usage

### StatsD
Initialize the `StatsDConfig` struct and call the function `InitiateStatsDMetrics` to initialize the client.
To use the statsDClient call the function `StatsDClient`.

### Sentry
Initialize the `SentryConfig` struct and call the function `InitSentry` to initialize the client.
To capture events use the following functions -
* `CaptureError`
* `CaptureWarn`
* `CaptureErrorWithTags`

### New Relic
Initialize the `newrelic.Config` struct and call the function `InitNewrelic` to initialize the client.
To capture events use the following functions -
* `StartPostgresDataSegmentNow`
* `StartMongoDBDataSegmentNow`
* `StartRedisSegmentNow`
* `StartSegmentNow`
* `StartKafkaPushSegment`
* `StartRabbitmqPushSegment`
* `StartExternalSegmentNow`

To create new contextual event tracing use following functions- 
* `NewHttpContext`
* `NewContextWithTransaction`
* `GetTx`