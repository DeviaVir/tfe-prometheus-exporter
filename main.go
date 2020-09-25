package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var TfeToken, TfeTokenPath, TfeAddress, VaultReadyPath string

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

func main() {

	listendAddr := getEnvDefault("HTTP_LISTENADDR", ":9112")

	TfeToken = getEnv("TFE_TOKEN")
	TfeTokenPath = getEnv("TFE_TOKEN_PATH")
	TfeAddress = getEnv("TFE_ADDRESS")
	VaultReadyPath = getEnv("VAULT_READY_PATH")

	if VaultReadyPath != "" {
		for {
			if fileExists(VaultReadyPath) {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	if TfeTokenPath != "" {
		if fileExists(TfeTokenPath) {
			path, err := homedir.Expand(TfeTokenPath)
			if err != nil {
				log.Fatal(err)
			}
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}

			TfeToken = strings.TrimSpace(string(content))
		}
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	tfeRuns := NewTfeCollector()
	prometheus.MustRegister(tfeRuns)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	log.Info("Now listening on ", listendAddr)
	log.Fatal(http.ListenAndServe(listendAddr, nil))

}
