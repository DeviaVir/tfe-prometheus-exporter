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
			[]string{"workspace"}, nil,

		),
		runsPendingMetric: prometheus.NewDesc("runs_pending",
			"Runs currently in the queue (pending)",
			[]string{"workspace"}, nil,
		),
		runsPlanningMetric: prometheus.NewDesc("runs_planning",
			"Runs currently planning (planning)",
			[]string{"workspace"}, nil,
		),
		runsPlannedMetric: prometheus.NewDesc("runs_planned",
			"Runs planned (planned)",
			[]string{"workspace"}, nil,
		),
		runsConfirmedMetric: prometheus.NewDesc("runs_confirmed",
			"Runs confirmed (confirmed)",
			[]string{"workspace"}, nil,
		),
		runsApplyingMetric: prometheus.NewDesc("runs_applying",
			"Runs currently applying (applying)",
			[]string{"workspace"}, nil,
		),
		runsAppliedMetric: prometheus.NewDesc("runs_applied",
			"Runs applied (applied)",
			[]string{"workspace"}, nil,
		),
		runsDiscardedMetric: prometheus.NewDesc("runs_discarded",
			"Runs discarded (discarded)",
			[]string{"workspace"}, nil,
		),
		runsErroredMetric: prometheus.NewDesc("runs_errored",
			"Runs errored (errored)",
			[]string{"workspace"}, nil,
		),
		runsCanceledMetric: prometheus.NewDesc("runs_canceled",
			"Runs canceled (canceled)",
			[]string{"workspace"}, nil,
		),
		runsPolicyCheckingMetric: prometheus.NewDesc("runs_policy_checking",
			"Runs currently checking policy (policy-checking)",
			[]string{"workspace"}, nil,
		),
		runsPolicyOverrideMetric: prometheus.NewDesc("runs_policy_override",
			"Runs with overriden policy (policy-override)",
			[]string{"workspace"}, nil,
		),
		runsPolicyCheckedMetric: prometheus.NewDesc("runs_policy_checked",
			"Runs with checked policy (policy-checked)",
			[]string{"workspace"}, nil,
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
func getWorkspaceNames(client *tfe.Client, ctx *context.Context, orgName string) (*[]string, error) {
	options := tfe.ListOptions{
		PageNumber: 1,
		PageSize:   1,
	}

	workspaces, err := client.Workspaces.List(*ctx, orgName, tfe.WorkspaceListOptions{ListOptions: options})
	if err != nil {
		return nil, err
	}
	var workspaceNames []string

	for _, workspace := range workspaces.Items {
		workspaceNames = append(workspaceNames, workspace.Name)
	}
	return &workspaceNames, err
}

// iterate over the list of workspace names, getting the matching runs for each
func getRunsByWorkspace(client *tfe.Client, ctx *context.Context, orgName string) (map[string]*tfe.AdminRunsList, error) {

	options := tfe.ListOptions{
		PageNumber: 1,
		PageSize:   1,
	}

	workspaces, err := getWorkspaceNames(client, ctx, orgName)
	if err != nil {
		return nil, err
	}

	runsPerWorkspace := make(map[string]*tfe.AdminRunsList)

	for _, workspace := range *workspaces {
		runs, err := client.AdminRuns.List(
			*ctx, tfe.AdminRunsListOptions{ListOptions: options, Query: &workspace})
		if err != nil {
			return nil, err
		} else {
			runsPerWorkspace[workspace] = runs
		}
	}

	return runsPerWorkspace, nil
}

//Collect implements required collect function for all prometheus collectors
func (collector *tfeCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[INFO]: scraping metrics per workspace")

	config := &tfe.Config{
		Token:   TfeToken,
		Address: TfeAddress,
	}
	ctx := context.Background()
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	runsPerWorkspace, err := getRunsByWorkspace(client, &ctx, TfeOrgName)
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
