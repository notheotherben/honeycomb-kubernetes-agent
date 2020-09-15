// Originally inspired by OpenTelemetry Collector kubeletstats receiver
// https://github.com/open-telemetry/opentelemetry-collector

package service

import (
	"github.com/honeycombio/honeycomb-kubernetes-agent/metrics"
	"github.com/honeycombio/libhoney-go"
	"github.com/sirupsen/logrus"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/honeycombio/honeycomb-kubernetes-agent/interval"
	"github.com/honeycombio/honeycomb-kubernetes-agent/kubelet"
)

var _ interval.Runnable = (*runnable)(nil)

type runnable struct {
	statsProvider         *kubelet.StatsProvider
	metadataProvider      *kubelet.MetadataProvider
	metricsProvider       *metrics.Processor
	interval              time.Duration
	k8sClusterName        string
	restClient            kubelet.RestClient
	builder               *libhoney.Builder
	omitLabels            []metrics.OmitLabel
	metricGroupsToCollect map[metrics.MetricGroup]bool
	logger                *logrus.Logger
}

func newRunnable(rc kubelet.RestClient, builder *libhoney.Builder, opt Options, logger *logrus.Logger) *runnable {
	return &runnable{
		interval:              opt.Interval,
		k8sClusterName:        opt.ClusterName,
		restClient:            rc,
		builder:               builder,
		omitLabels:            opt.OmitLabels,
		metricGroupsToCollect: opt.MetricGroupsToCollect,
		logger:                logger,
	}
}

// Sets up the kubelet connection at startup time.
func (r *runnable) Setup() error {
	r.logger.Info("Creating Metrics Service Providers...")

	r.statsProvider = kubelet.NewStatsProvider(r.restClient)
	r.metadataProvider = kubelet.NewMetadataProvider(r.restClient)
	r.metricsProvider = metrics.NewMetricsProcessor(r.interval, r.logger)

	r.logger.Debug("Metrics Service Providers created")
	return nil
}

func (r *runnable) Run() error {

	// get all stats from stats provider
	summary, err := r.statsProvider.StatsSummary()
	if err != nil {
		logrus.WithError(err).Error("Could not retrieve stats")
		return nil
	}
	r.logger.WithFields(logrus.Fields{
		"podCount": len(summary.Pods),
	}).Debug("Retrieved Stats")

	// fetch metadata for all pods/containers
	var podsMetadata *v1.PodList
	podsMetadata, err = r.metadataProvider.Pods()
	if err != nil {
		r.logger.WithError(err).Error("Could not retrieve metadata")
		return nil
	}
	r.logger.WithFields(logrus.Fields{
		"podCount": len(podsMetadata.Items),
	}).Debug("Retrieved Metadata")

	// Get resource metrics
	metadata := metrics.NewMetadata(podsMetadata, r.omitLabels, r.logger)
	resourceMetrics := r.metricsProvider.GenerateMetricsData(summary, metadata, r.metricGroupsToCollect)
	r.logger.WithFields(logrus.Fields{
		"resourceCount": len(resourceMetrics),
	}).Debug("Processing Metrics Data...")

	// iterate over resource metrics data
	for _, rm := range resourceMetrics {

		ev, err := r.createEventFromResource(rm.Resource)
		if err != nil {
			r.logger.WithFields(logrus.Fields{
				"resourceType": rm.Resource.Type,
				"resourceName": rm.Resource.Name,
			}).WithError(err).Error("Could not create event for resource")
			continue
		}

		r.logger.WithFields(logrus.Fields{
			"resourceType": rm.Resource.Type,
			"resourceName": rm.Resource.Name,
		}).Trace("Creating event for resource")

		pre := metrics.PrefixMetrics

		for k, v := range rm.Metrics {
			r.logger.WithFields(logrus.Fields{
				"resourceType": rm.Resource.Type,
				"resourceName": rm.Resource.Name,
				"metricName":   k,
			}).Trace("Metric to event field")

			var val float64
			if v.IsCounter {
				val = r.metricsProvider.GetCounterRate(rm.Resource, k, v)
			} else {
				val = v.GetValue()
			}
			ev.AddField(pre+k, val)
		}

		r.logger.WithFields(logrus.Fields{
			"resourceType": rm.Resource.Type,
			"resourceName": rm.Resource.Name,
			"metricCount":  len(rm.Metrics),
		}).Debug("Event Data: ", ev.Fields())

		if err = ev.Send(); err != nil {
			r.logger.WithFields(logrus.Fields{
				"resourceType": rm.Resource.Type,
				"resourceName": rm.Resource.Name,
				"metricCount":  len(rm.Metrics),
			}).WithError(err).Error("Error sending event to honeycomb")
		}
	}

	return nil
}

func (r *runnable) createEventFromResource(res *metrics.Resource) (*libhoney.Event, error) {
	ev := r.builder.NewEvent()
	_ = ev.Add(map[string]string{
		metrics.PrefixMetrics + metrics.MetricSourceName: res.Name,
		metrics.PrefixMetrics + metrics.MetricSourceType: res.Type,
		metrics.KubernetesResourceType:                   res.Type,
		metrics.PrefixCluster + "name":                   r.k8sClusterName,
	})
	ev.Timestamp = res.Timestamp
	if err := ev.Add(res.Labels); err != nil {
		return nil, err
	}
	if res.Status != nil {
		if err := ev.Add(res.Status); err != nil {
			return nil, err
		}
	}
	return ev, nil
}
