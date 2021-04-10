package metrics

import (
	"context"
	"log"
	"net/http"

	"time"

	newrelic "github.com/newrelic/go-agent"
)

type nCtxKey int

const txKey nCtxKey = 0

var newrelicApp newrelic.Application

func InitNewrelic(cfg newrelic.Config) {
	if cfg.Enabled {
		var err error
		newrelicApp, err = newrelic.NewApplication(cfg)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func ShutdownNewrelic() {
	if newrelicApp != nil {
		newrelicApp.Shutdown(time.Second)
	}
}

func NewrelicApp() newrelic.Application {
	return newrelicApp
}

func StartPostgresDataSegmentNow(op string, tableName string, txn newrelic.Transaction) *newrelic.DatastoreSegment {
	return startDatastoreSegment(op, tableName, txn, newrelic.DatastorePostgres)
}

func StartMongoDBDataSegmentNow(op string, collectionName string, txn newrelic.Transaction) *newrelic.DatastoreSegment {
	return startDatastoreSegment(op, collectionName, txn, newrelic.DatastoreMongoDB)
}

func StartCassandraSegmentNow(op string, tableName string, txn newrelic.Transaction) *newrelic.DatastoreSegment {
	return startDatastoreSegment(op, tableName, txn, newrelic.DatastoreCassandra)
}

func StartRedisSegmentNow(op string, tableName string, txn newrelic.Transaction) *newrelic.DatastoreSegment {
	return startDatastoreSegment(op, tableName, txn, newrelic.DatastoreRedis)
}

func startDatastoreSegment(op string, tableName string, txn newrelic.Transaction, product newrelic.DatastoreProduct) *newrelic.DatastoreSegment {
	s := newrelic.DatastoreSegment{
		Product:    product,
		Collection: tableName,
		Operation:  op,
		StartTime:  newrelic.StartSegmentNow(txn),
	}

	return &s
}

func StartSegmentNow(name string, txn newrelic.Transaction) *newrelic.Segment {
	s := newrelic.Segment{
		Name:      name,
		StartTime: newrelic.StartSegmentNow(txn),
	}

	return &s
}

func StartKafkaPushSegment(txn newrelic.Transaction, topic string) *newrelic.MessageProducerSegment {
	s := newrelic.MessageProducerSegment{
		StartTime:            newrelic.StartSegmentNow(txn),
		Library:              "Kafka",
		DestinationType:      newrelic.MessageTopic,
		DestinationName:      topic,
		DestinationTemporary: false,
	}

	return &s
}

func StartRabbitmqPushSegment(txn newrelic.Transaction, exchange string) *newrelic.MessageProducerSegment {
	s := newrelic.MessageProducerSegment{
		StartTime:            newrelic.StartSegmentNow(txn),
		Library:              "RabbitMQ",
		DestinationType:      newrelic.MessageExchange,
		DestinationName:      exchange,
		DestinationTemporary: false,
	}

	return &s
}

func StartExternalSegmentNow(txn newrelic.Transaction, url string) *newrelic.ExternalSegment {
	s := newrelic.ExternalSegment{
		StartTime: newrelic.StartSegmentNow(txn),
		URL:       url,
	}

	return &s
}

func NewHTTPContext(ctx context.Context, w http.ResponseWriter) context.Context {
	if newrelicApp != nil {
		tx, ok := w.(newrelic.Transaction)
		if !ok {
			return ctx
		}
		return context.WithValue(ctx, txKey, tx)
	}
	return ctx
}

func NewContextWithTransaction(ctx context.Context, tx newrelic.Transaction) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) (newrelic.Transaction, bool) {
	tx, ok := ctx.Value(txKey).(newrelic.Transaction)
	return tx, ok
}
