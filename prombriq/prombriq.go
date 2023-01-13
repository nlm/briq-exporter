package prombriq

import (
	"context"
	"time"

	"github.com/nlm/briq-cli/briq"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	briqActiveBalanceDesc   = prometheus.NewDesc("briq_active_balance", "active balance of a user", []string{"user"}, nil)
	briqInactiveBalanceDesc = prometheus.NewDesc("briq_inactive_balance", "inactive balance of a user", []string{"user"}, nil)
	briqPointsDesc          = prometheus.NewDesc("briq_points", "points earned by a user", []string{"user"}, nil)
)

// Option is an option for a Collector
type Option func(*Collector)

// NewBriqCollector creates a new Briq collector from a client
func NewCollector(client *briq.Client, options ...Option) *Collector {
	c := &Collector{client: client}
	for _, option := range options {
		option(c)
	}
	return c
}

// WithLogger sets the logger of a briq collector
func WithLogger(logger *log.Logger) Option {
	return func(c *Collector) {
		c.logger = logger
	}
}

// WithTimeout sets a timeout for briq api queries
func WithTimeout(timeout time.Duration) Option {
	return func(c *Collector) {
		c.timeout = timeout
	}
}

// Collector is a prometheus collector to collect briq metrics
type Collector struct {
	client  *briq.Client
	logger  *log.Logger
	timeout time.Duration
}

// Collect is part of the implementation of the prometheus collector interface
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	// Setup Context
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx = context.Background()
	if c.timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// Collect metrics
	res, err := c.client.ListUsers(ctx, &briq.ListUsersRequest{})
	if err != nil {
		c.logger.WithError(err).Error("error calling briq api")
		return
	}
	for _, user := range res.Users {
		ch <- prometheus.MustNewConstMetric(
			briqActiveBalanceDesc,
			prometheus.GaugeValue,
			float64(user.ActiveBalance),
			user.Username,
		)
		ch <- prometheus.MustNewConstMetric(
			briqInactiveBalanceDesc,
			prometheus.GaugeValue,
			float64(user.InactiveBalance),
			user.Username,
		)
		ch <- prometheus.MustNewConstMetric(
			briqPointsDesc,
			prometheus.GaugeValue,
			float64(user.Points),
			user.Username,
		)
	}
}

// Describe is part of the implementation of the prometheus collector interface
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- briqActiveBalanceDesc
	ch <- briqInactiveBalanceDesc
	ch <- briqPointsDesc
}
