package metrics

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// Database metrics
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type", "table"},
	)

	DBConnectionsActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
		[]string{"database"},
	)

	// Redis metrics
	RedisCommandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_commands_total",
			Help: "Total number of Redis commands",
		},
		[]string{"command", "status"},
	)

	RedisCommandDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_command_duration_seconds",
			Help:    "Redis command duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"command"},
	)

	// Kafka metrics
	KafkaMessagesProduced = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_produced_total",
			Help: "Total number of Kafka messages produced",
		},
		[]string{"topic"},
	)

	KafkaMessagesConsumed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total number of Kafka messages consumed",
		},
		[]string{"topic", "consumer_group"},
	)

	// Business metrics
	BookingAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_attempts_total",
			Help: "Total number of booking attempts",
		},
		[]string{"status"}, // success, failed, cancelled
	)

	PaymentsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payments_total",
			Help: "Total number of payments",
		},
		[]string{"status"}, // success, failed, pending
	)

	NotificationsSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notifications_sent_total",
			Help: "Total number of notifications sent",
		},
		[]string{"type"}, // email, sms, push
	)

	// Worker metrics
	WorkerJobsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_jobs_processed_total",
			Help: "Total number of worker jobs processed",
		},
		[]string{"job_type", "status"},
	)

	WorkerJobDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "worker_job_duration_seconds",
			Help:    "Worker job duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"job_type"},
	)
)

// GinMiddleware returns a Gin middleware that collects HTTP metrics
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		HTTPRequestsInFlight.Inc()

		// Start timer
		timer := prometheus.NewTimer(HTTPRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()))

		c.Next()

		// Stop timer
		timer.ObserveDuration()

		// Record request
		HTTPRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(c.Writer.Status())).Inc()
		HTTPRequestsInFlight.Dec()
	}
}

// RecordDBQuery records database query duration
func RecordDBQuery(queryType, table string, duration float64) {
	DBQueryDuration.WithLabelValues(queryType, table).Observe(duration)
}

// RecordRedisCommand records Redis command execution
func RecordRedisCommand(command, status string, duration float64) {
	RedisCommandsTotal.WithLabelValues(command, status).Inc()
	RedisCommandDuration.WithLabelValues(command).Observe(duration)
}

// RecordKafkaMessage records Kafka message production/consumption
func RecordKafkaMessageProduced(topic string) {
	KafkaMessagesProduced.WithLabelValues(topic).Inc()
}

func RecordKafkaMessageConsumed(topic, consumerGroup string) {
	KafkaMessagesConsumed.WithLabelValues(topic, consumerGroup).Inc()
}

// RecordBookingAttempt records booking attempt
func RecordBookingAttempt(status string) {
	BookingAttemptsTotal.WithLabelValues(status).Inc()
}

// RecordPayment records payment
func RecordPayment(status string) {
	PaymentsTotal.WithLabelValues(status).Inc()
}

// RecordNotification records notification sent
func RecordNotification(notificationType string) {
	NotificationsSent.WithLabelValues(notificationType).Inc()
}

// RecordWorkerJob records worker job processing
func RecordWorkerJob(jobType, status string, duration float64) {
	WorkerJobsProcessed.WithLabelValues(jobType, status).Inc()
	WorkerJobDuration.WithLabelValues(jobType).Observe(duration)
}

