package main

import (
	"context"

	tfe "github.com/DeviaVir/go-tfe"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type tfeCollector struct {
	runsTotalMetric          *prometheus.Desc
	runsPendingMetric        *prometheus.Desc
	runsPlanningMetric       *prometheus.Desc
	runsPlannedMetric        *prometheus.Desc
	runsConfirmedMetric      *prometheus.Desc
	runsApplyingMetric       *prometheus.Desc
	runsAppliedMetric        *prometheus.Desc
	runsDiscardedMetric      *prometheus.Desc
	runsErroredMetric        *prometheus.Desc
	runsCanceledMetric       *prometheus.Desc
	runsPolicyCheckingMetric *prometheus.Desc
	runsPolicyOverrideMetric *prometheus.Desc
	runsPolicyCheckedMetric  *prometheus.Desc
}

// newTFECollector initializes every descriptor and returns a pointer to the collector
func NewTfeCollector() *tfeCollector {
	return &tfeCollector{
		runsTotalMetric: prometheus.NewDesc("runs_total",
			"Total number of runs with any status (total)",
			nil, nil,
		),
		runsPendingMetric: prometheus.NewDesc("runs_pending",
			"Runs currently in the queue (pending)",
			nil, nil,
		),
		runsPlanningMetric: prometheus.NewDesc("runs_planning",
			"Runs currently planning (planning)",
			nil, nil,
		),
		runsPlannedMetric: prometheus.NewDesc("runs_planned",
			"Runs planned (planned)",
			nil, nil,
		),
		runsConfirmedMetric: prometheus.NewDesc("runs_confirmed",
			"Runs confirmed (confirmed)",
			nil, nil,
		),
		runsApplyingMetric: prometheus.NewDesc("runs_applying",
			"Runs currently applying (applying)",
			nil, nil,
		),
		runsAppliedMetric: prometheus.NewDesc("runs_applied",
			"Runs applied (applied)",
			nil, nil,
		),
		runsDiscardedMetric: prometheus.NewDesc("runs_discarded",
			"Runs discarded (discarded)",
			nil, nil,
		),
		runsErroredMetric: prometheus.NewDesc("runs_errored",
			"Runs errored (errored)",
			nil, nil,
		),
		runsCanceledMetric: prometheus.NewDesc("runs_canceled",
			"Runs canceled (canceled)",
			nil, nil,
		),
		runsPolicyCheckingMetric: prometheus.NewDesc("runs_policy_checking",
			"Runs currently checking policy (policy-checking)",
			nil, nil,
		),
		runsPolicyOverrideMetric: prometheus.NewDesc("runs_policy_override",
			"Runs with overriden policy (policy-override)",
			nil, nil,
		),
		runsPolicyCheckedMetric: prometheus.NewDesc("runs_policy_checked",
			"Runs with checked policy (policy-checked)",
			nil, nil,
		),
	}
}

// Describe essentially writes all descriptors to the prometheus desc channel.
func (collector *tfeCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.runsTotalMetric
	ch <- collector.runsPendingMetric
	ch <- collector.runsPlanningMetric
	ch <- collector.runsPlannedMetric
	ch <- collector.runsConfirmedMetric
	ch <- collector.runsApplyingMetric
	ch <- collector.runsAppliedMetric
	ch <- collector.runsDiscardedMetric
	ch <- collector.runsErroredMetric
	ch <- collector.runsCanceledMetric
	ch <- collector.runsPolicyCheckingMetric
	ch <- collector.runsPolicyOverrideMetric
	ch <- collector.runsPolicyCheckedMetric
}

// tfeClients returns a List all the runs of the terraform enterprise installation
func tfeRuns(tfeToken, tfeAddress string) (*tfe.AdminRunsList, error) {

	config := &tfe.Config{
		Token:   tfeToken,
		Address: tfeAddress,
	}
	ctx := context.Background()
	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}

	options := tfe.ListOptions{
		PageNumber: 1,
		PageSize:   1,
	}
	runs, err := client.AdminRuns.List(
		ctx, tfe.AdminRunsListOptions{ListOptions: options})
	if err != nil {
		return nil, err
	}

	return runs, nil
}

//Collect implements required collect function for all promehteus collectors
func (collector *tfeCollector) Collect(ch chan<- prometheus.Metric) {

	log.Println("[INFO]: scraping metrics")

	// runs retrieves a List all the runs of TFE.
	runs, err := tfeRuns(TfeToken, TfeAddress)
	if err != nil {
		log.Fatal(err)
	}

	//Write latest value for each metric in the prometheus metric channel.
	ch <- prometheus.MustNewConstMetric(collector.runsTotalMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Total))
	ch <- prometheus.MustNewConstMetric(collector.runsPendingMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Pending))
	ch <- prometheus.MustNewConstMetric(collector.runsPlanningMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Planning))
	ch <- prometheus.MustNewConstMetric(collector.runsPlannedMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Planned))
	ch <- prometheus.MustNewConstMetric(collector.runsConfirmedMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Confirmed))
	ch <- prometheus.MustNewConstMetric(collector.runsApplyingMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Applying))
	ch <- prometheus.MustNewConstMetric(collector.runsAppliedMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Applied))
	ch <- prometheus.MustNewConstMetric(collector.runsDiscardedMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Discarded))
	ch <- prometheus.MustNewConstMetric(collector.runsErroredMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Errored))
	ch <- prometheus.MustNewConstMetric(collector.runsCanceledMetric, prometheus.GaugeValue, float64(runs.StatusCounts.Canceled))
	ch <- prometheus.MustNewConstMetric(collector.runsPolicyCheckingMetric, prometheus.GaugeValue, float64(runs.StatusCounts.PolicyChecking))
	ch <- prometheus.MustNewConstMetric(collector.runsPolicyOverrideMetric, prometheus.GaugeValue, float64(runs.StatusCounts.PolicyOverride))
	ch <- prometheus.MustNewConstMetric(collector.runsPolicyCheckedMetric, prometheus.GaugeValue, float64(runs.StatusCounts.PolicyChecked))
}
