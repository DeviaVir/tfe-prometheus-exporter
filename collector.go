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
func NewTfeCollector(tfeToken, tfeAddress string) *tfeCollector {
	workspaces, err := getWorkspaceNames(tfeToken, tfeAddress)
	if err != nil {
		log.Fatal(err)
	}

	return &tfeCollector{
		runsTotalMetric: prometheus.NewDesc("runs_total",
			"Total number of runs with any status (total)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsPendingMetric: prometheus.NewDesc("runs_pending",
			"Runs currently in the queue (pending)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsPlanningMetric: prometheus.NewDesc("runs_planning",
			"Runs currently planning (planning)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsPlannedMetric: prometheus.NewDesc("runs_planned",
			"Runs planned (planned)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsConfirmedMetric: prometheus.NewDesc("runs_confirmed",
			"Runs confirmed (confirmed)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsApplyingMetric: prometheus.NewDesc("runs_applying",
			"Runs currently applying (applying)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsAppliedMetric: prometheus.NewDesc("runs_applied",
			"Runs applied (applied)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsDiscardedMetric: prometheus.NewDesc("runs_discarded",
			"Runs discarded (discarded)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsErroredMetric: prometheus.NewDesc("runs_errored",
			"Runs errored (errored)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsCanceledMetric: prometheus.NewDesc("runs_canceled",
			"Runs canceled (canceled)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsPolicyCheckingMetric: prometheus.NewDesc("runs_policy_checking",
			"Runs currently checking policy (policy-checking)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsPolicyOverrideMetric: prometheus.NewDesc("runs_policy_override",
			"Runs with overriden policy (policy-override)",
			nil, prometheus.Labels{"workspace": workspaces},
		),
		runsPolicyCheckedMetric: prometheus.NewDesc("runs_policy_checked",
			"Runs with checked policy (policy-checked)",
			nil, prometheus.Labels{"workspace": workspaces},
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


// get a list of workspace names from TFE
func getWorkspaceNames(tfeToken, tfeAddress string) ([]string, error) {
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

	workspaces, err := client.Workspaces.List(ctx, tfe.WorkspaceListOptions{ListOptions: options})
	if err != nil {
		return nil, err
	}
	var workspaceNames []string

	for workspace := range workspaces {
		workspaceNames = append(workspaceNames, workspace.Name)
	}
	return workspaceNames, err
}

// iterate over the list of workspace names, getting the matching runs for each
func getRunsByWorkspace(tfeToken, tfeAddress string) (map[string]*tfe.AdminRunsList, error) {
	config := &tfe.Config{
		Token:   tfeToken,
		Address: tfeAddress,
	}
	ctx := context.Background()
	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}

	workspaces, err := getWorkspaceNames(tfeToken, tfeAddress)
	if err != nil {
		return nil, err
	}

	runsPerWorkspace := make(map[string]*tfe.AdminRunsList)

	for workspace := range workspaces {
		runs, err := client.AdminRuns.List(
			// this is just a string match; it's probably going to choke with our service manager naming construct?
			ctx, tfe.AdminRunsListOptions{ListOptions: options, Query: workspace.Name})
		if err != nil {
			return nil, err
		} else {
			runsPerWorkspace[workspace.Name] = runs
		}
	}

	return runsPerWorkspace, nil
}

//Collect implements required collect function for all prometheus collectors
func (collector *tfeCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[INFO]: scraping metrics per workspace")

	runsPerWorkspace, err := getRunsByWorkspace(TfeToken, TfeAddress)
	if err != nil {
		log.Fatal(err)
	}

	for workspaceName, runList := range runsPerWorkspace {
		//Write latest value for each metric in the prometheus metric channel.
		ch <- prometheus.MustNewConstMetric(collector.runsTotalMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Total), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsPendingMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Pending), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsPlanningMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Planning), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsPlannedMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Planned), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsConfirmedMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Confirmed), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsApplyingMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Applying), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsAppliedMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Applied), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsDiscardedMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Discarded), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsErroredMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Errored), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsCanceledMetric, prometheus.GaugeValue, float64(runList.StatusCounts.Canceled), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsPolicyCheckingMetric, prometheus.GaugeValue, float64(runList.StatusCounts.PolicyChecking), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsPolicyOverrideMetric, prometheus.GaugeValue, float64(runList.StatusCounts.PolicyOverride), workspaceName)
		ch <- prometheus.MustNewConstMetric(collector.runsPolicyCheckedMetric, prometheus.GaugeValue, float64(runList.StatusCounts.PolicyChecked), workspaceName)
	}
}
