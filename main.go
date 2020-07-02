package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	tfe "github.com/DeviaVir/go-tfe"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getEnv(name string) string {
	envValue, ok := os.LookupEnv(name)
	if ok {
		return envValue
	}
	return ""
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
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	tfeToken := getEnv("TFE_TOKEN")
	tfeTokenPath := getEnv("TFE_TOKEN_PATH")
	tfeAddress := getEnv("TFE_ADDRESS")
	listendAddr := getEnvDefault("HTTP_LISTENADDR", ":9112")

	if tfeTokenPath != "" {
		if fileExists(tfeTokenPath) {
			path, err := homedir.Expand(tfeTokenPath)
			if err != nil {
				log.Fatal(err)
			}
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}

			tfeToken = strings.TrimSpace(string(content))
		}
	}

	config := &tfe.Config{
		Token:   tfeToken,
		Address: tfeAddress,
	}
	ctx := context.Background()
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	options := tfe.ListOptions{
		PageNumber: 0,
		PageSize:   0,
	}
	runs, err := client.AdminRuns.List(
		ctx, tfe.AdminRunsListOptions{ListOptions: options})
	if err != nil {
		log.Fatal(err)
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
	log.Info("Now listening on ", listendAddr)
	log.Fatal(http.ListenAndServe(listendAddr, nil))
}
