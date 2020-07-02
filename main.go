package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/DeviaVir/go-tfe"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func getEnv(name string) string {
	envValue, ok := os.LookupEnv(name)
	if ok {
		return envValue
	}
	panic(fmt.Sprintf("Missing environment variable: %s", name))
}

func getEnvDefault(name string, defaultVal string) string {
	envValue, ok := os.LookupEnv(name)
	if ok {
		return envValue
	}
	return defaultVal
}

func setGauge(name string, help string, callback func() float64) {
	gaugeFunc := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "tf",
		Subsystem: "enterprise",
		Name:      name,
		Help:      help,
	}, callback)
	prometheus.MustRegister(gaugeFunc)
}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	tfeToken := getEnv("TFE_TOKEN")
	tfeAddress := getEnv("TFE_ADDRESS")
	listendAddr := getEnvDefault("HTTP_LISTENADDR", ":9112")

	config := &tfe.Config{
		Token:   tfeToken,
		Address: tfeAddress,
	}
	ctx = context.Background()
	client, err := tfe.NewClient(config)
	if err != nil {
		panic(err)
	}

	options := tfe.ListOptions{
		PageNumber: 0,
		PageSize:   0,
	}
	runs, err := client.AdminRuns.List(
		ctx, tfe.AdminRunsListOptions{ListOptions: options})
	if err != nil {
		panic(err)
	}

	setGauge("runs_total", "Total number of runs with any status (total)", func() float64 {
		return float64(runs.StatusCounts.Total)
	})
	setGauge("runs_pending", "Runs currently in the queue (pending)", func() float64 {
		return float64(runs.StatusCounts.Pending)
	})
	setGauge("runs_planning", "Runs currently planning (planning)", func() float64 {
		return float64(runs.StatusCounts.Planning)
	})
	setGauge("runs_planned", "Runs planned (planned)", func() float64 {
		return float64(runs.StatusCounts.Planned)
	})
	setGauge("runs_confirmed", "Runs confirmed (confirmed)", func() float64 {
		return float64(runs.StatusCounts.Confirmed)
	})
	setGauge("runs_applying", "Runs currently applying (applying)", func() float64 {
		return float64(runs.StatusCounts.Applying)
	})
	setGauge("runs_applied", "Runs applied (applied)", func() float64 {
		return float64(runs.StatusCounts.Applied)
	})
	setGauge("runs_discarded", "Runs discarded (discarded)", func() float64 {
		return float64(runs.StatusCounts.Discarded)
	})
	setGauge("runs_errored", "Runs errored (errored)", func() float64 {
		return float64(runs.StatusCounts.Errored)
	})
	setGauge("runs_canceled", "Runs canceled (canceled)", func() float64 {
		return float64(runs.StatusCounts.Canceled)
	})
	setGauge("runs_policy_checking", "Runs currently checking policy (policy-checking)", func() float64 {
		return float64(runs.StatusCounts.PolicyChecking)
	})
	setGauge("runs_policy_override", "Runs with overriden policy (policy-override)", func() float64 {
		return float64(runs.StatusCounts.PolicyOverride)
	})
	setGauge("runs_policy_checked", "Runs with checked policy (policy-checked)", func() float64 {
		return float64(runs.StatusCounts.PolicyChecked)
	})

	http.Handle("/metrics", promhttp.Handler())
	logrus.Info("Now listening on ", listendAddr)
	logrus.Fatal(http.ListenAndServe(listendAddr, nil))
}
