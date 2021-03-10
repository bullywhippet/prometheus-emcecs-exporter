package collector

import (
	"github.com/paychex/prometheus-emcecs-exporter/pkg/ecsclient"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// A EcsNodeSystemStatsCollector implements the prometheus.Collector.
type EcsNodeSystemStatsCollector struct {
	ecsClient *ecsclient.EcsClient
	namespace string
}

var (
	nodeCpuUtilization = prometheus.NewDesc(
		prometheus.BuildFQName("emcecs", "node", "cpuUtilizationPercent"),
		"Average current CPU utilization percent on node",
		[]string{"node"}, nil,
	)
	nodeMemoryUtilization = prometheus.NewDesc(
		prometheus.BuildFQName("emcecs", "node", "memoryUtilizationPercent"),
		"Average current memory utilization percent on node",
		[]string{"node"}, nil,
	)
	activeConnections = prometheus.NewDesc(
		prometheus.BuildFQName("emcecs", "node", "activeConnections"),
		"Number of current active connections on node",
		[]string{"node"}, nil,
	)
)

// NewNodeSystemStatsCollector returns an initialized Node DT Collector.
func NewNodeSystemStatsCollector(emcecs *ecsclient.EcsClient, namespace string) (*EcsNodeSystemStatsCollector, error) {

	log.WithFields(log.Fields{"package": "node-system-stats-collector"}).Debug("Init Node exporter")
	return &EcsNodeSystemStatsCollector{
		ecsClient: emcecs,
		namespace: namespace,
	}, nil
}

// Collect fetches the stats from configured nodes as Prometheus metrics.
// It implements prometheus.Collector.
func (e *EcsNodeSystemStatsCollector) Collect(ch chan<- prometheus.Metric) {
	log.WithFields(log.Fields{"package": "node-system-stats-collector"}).Debug("ECS Node DT collect starting")
	if e.ecsClient == nil {
		log.WithFields(log.Fields{"package": "node-system-stats-collector"}).Error("ECS client not configured.")
		return
	}

	nodeSystemStats := e.ecsClient.RetrieveNodeystemStatusParallel()
	for _, node := range nodeSystemStats {
		ch <- prometheus.MustNewConstMetric(nodeCpuUtilization, prometheus.GaugeValue, node.CPUUtilization, node.NodeIP)
		ch <- prometheus.MustNewConstMetric(nodeMemoryUtilization, prometheus.GaugeValue, node.MemoryUtilization, node.NodeIP)

		ch <- prometheus.MustNewConstMetric(activeConnections, prometheus.GaugeValue, node.ActiveConnections, node.NodeIP)
	}

	log.WithFields(log.Fields{"package": "node-system-stats-collector"}).Debug("NodeSystemStats exporter finished")
	log.WithFields(log.Fields{"package": "node-system-stats-collector"}).Debug(nodeSystemStats)
}

// Describe describes the metrics exported from this collector.
func (e *EcsNodeSystemStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nodeCpuUtilization
	ch <- nodeMemoryUtilization
	ch <- activeConnections
}
