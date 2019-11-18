// Copyright 2018 Andrea Spagnolo
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	//"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	//"os"
	"os/exec"
	//"strconv"
	//"strings"
	//"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9367").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	dhcpdPidFile  = kingpin.Flag("dhcpd.pid-file", "Path where dhcpd PID file is located.").Default("/var/run/dhcpd.pid").String()
	dhcpdConfigFile  = kingpin.Flag("dhcpd.config-file", "Path where dhcpd main config file (dhcpd.conf) is located.").Default("/etc/dhcp/dhcpd.conf").String()
	dhcpdLeaseFile  = kingpin.Flag("dhcpd.lease-file", "Path where dhcpd lease file is located.").Default("/var/lib/dhcp/dhcpd.leases").String()
)

const (
	namespace = "iscdhcpd"
)

type Exporter struct {
	summaryFree          *prometheus.Desc
	summaryTouched       *prometheus.Desc
	summaryUsed          *prometheus.Desc
	summaryDefined       *prometheus.Desc
	subnetFree           *prometheus.Desc
	subnetTouched        *prometheus.Desc
	subnetUsed           *prometheus.Desc
	subnetDefined        *prometheus.Desc
	sharedNetworkFree    *prometheus.Desc
	sharedNetworkTouched *prometheus.Desc
	sharedNetworkUsed    *prometheus.Desc
	sharedNetworkDefined *prometheus.Desc
	//dhcpdUp              *prometheus.Desc
	scrapeFailures       prometheus.Counter
}

type Subnet struct {
	Location string  `json:"location"`
	Range    string  `json:"range"`
	Defined  float64 `json:"defined"`
	Used     float64 `json:"used"`
	Touched  float64 `json:"touched"`
	Free     float64 `json:"free"`
}

type SharedNetwork struct {
	Location string  `json:"location"`
	Defined  float64 `json:"defined"`
	Used     float64 `json:"used"`
	Touched  float64 `json:"touched"`
	Free     float64 `json:"free"`
}

type PoolsStats struct {
	Subnets        []Subnet        `json:"subnets"`
	SharedNetworks []SharedNetwork `json:"shared-networks"`
	Summary        Subnet          `json:"summary"`
}

func NewExporter() *Exporter {
	return &Exporter{
		summaryFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "summary", "free"),
			"Overall IPs Free",
			nil,
			nil),
		summaryTouched: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "summary", "touched"),
			"Overall IPs Touched",
			nil,
			nil),
		summaryUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "summary", "used"),
			"Overall IPs Used",
			nil,
			nil),
		summaryDefined: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "summary", "defined"),
			"Overall IPs Defined",
			nil,
			nil),
		subnetFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "subnet", "free"),
			"Subnet IPs Defined", []string{"location", "range"}, nil,
		),
		subnetTouched: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "subnet", "touched"),
			"Subnet IPs Touched", []string{"location", "range"}, nil,
		),
		subnetUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "subnet", "used"),
			"Subnet IPs Used", []string{"location", "range"}, nil,
		),
		subnetDefined: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "subnet", "defined"),
			"Subnet IPs Defined", []string{"location", "range"}, nil,
		),
		sharedNetworkFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "shared_network", "free"),
			"Shared Network IPs Defined", []string{"location"}, nil,
		),
		sharedNetworkTouched: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "shared_network", "touched"),
			"Shared Network IPs Touched", []string{"location"}, nil,
		),
		sharedNetworkUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "shared_network", "used"),
			"Shared Network IPs Used", []string{"location"}, nil,
		),
		sharedNetworkDefined: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "shared_network", "defined"),
			"Shared Network IPs Defined", []string{"location"}, nil,
		),
		// dhcpdUp: prometheus.NewDesc(
		// 	prometheus.BuildFQName(namespace, "process", "up"),
		// 	"Whether dhcpd daemon is running at PID defined at its pid-file.",
		// 	[]string{}, nil,
		// ),
		scrapeFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrape_failures",
			Help:      "Number of errors while scraping dhcpd-pools.",
		}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.summaryFree
	ch <- e.summaryTouched
	ch <- e.summaryUsed
	ch <- e.summaryDefined

	ch <- e.subnetFree
	ch <- e.subnetTouched
	ch <- e.subnetUsed
	ch <- e.subnetDefined

	ch <- e.sharedNetworkFree
	ch <- e.sharedNetworkTouched
	ch <- e.sharedNetworkUsed
	ch <- e.sharedNetworkDefined
	//ch <- e.dhcpdUp

	e.scrapeFailures.Describe(ch)
}

func getoutputPool() PoolsStats {
	outputPools, err := exec.Command("/usr/bin/dhcpd-pools", "-c", *dhcpdConfigFile, "--leases=" + *dhcpdLeaseFile, "-f", "j").Output()
	if err != nil {
		log.Errorf("Error: %s", err)
	}

	var poolsStats PoolsStats
	err = json.Unmarshal(outputPools, &poolsStats)

	if err != nil {
		log.Errorf("Error: %s", err)
	}

	return poolsStats
}

// func getDhcpdStatus() (status int) {
// 	status = 0
// 	strPid, err := ioutil.ReadFile(*dhcpdPidFile)
// 	if err != nil {
// 		log.Debugf("Unable to read DHCPD PID file: %s", err)
// 		return
// 	}

// 	pid, err := strconv.Atoi(strings.Trim(string(strPid), "\n"))
// 	if err != nil {
// 		log.Debugf("Unable to convert DHCP PID to integer: %s", err)
// 		return
// 	}
// 	log.Debugf("Read PID from dhcpd pid file: %d", pid)

// 	proc, err := os.FindProcess(pid)
// 	if err != nil {
// 		log.Debugf("Unknown error when retrieving process %d info: %s", pid, err)
// 		return
// 	}

// 	if err = proc.Signal(syscall.Signal(0)); err != nil {
// 		log.Debugf("DHCP process at PID %d is not running: %s", pid, err)
// 		return
// 	}
// 	status = 1
// 	log.Debugf("DHCP process with PID %d is running", pid)
// 	return
// }

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	// ch <- prometheus.MustNewConstMetric(
	// 	e.dhcpdUp, prometheus.GaugeValue, float64(0), //getDhcpdStatus()),
	// )
	outputPool := getoutputPool()

	ch <- prometheus.MustNewConstMetric(
		e.summaryFree, prometheus.GaugeValue, outputPool.Summary.Free,
	)
	ch <- prometheus.MustNewConstMetric(
		e.summaryTouched, prometheus.GaugeValue, outputPool.Summary.Touched,
	)
	ch <- prometheus.MustNewConstMetric(
		e.summaryUsed, prometheus.GaugeValue, outputPool.Summary.Used,
	)
	ch <- prometheus.MustNewConstMetric(
		e.summaryDefined, prometheus.GaugeValue, outputPool.Summary.Defined,
	)

	for subnet := 0; subnet < len(outputPool.Subnets); subnet++ {
		ch <- prometheus.MustNewConstMetric(
			e.subnetFree, prometheus.GaugeValue, outputPool.Subnets[subnet].Free,
			outputPool.Subnets[subnet].Location, outputPool.Subnets[subnet].Range,
		)
		ch <- prometheus.MustNewConstMetric(
			e.subnetTouched, prometheus.GaugeValue, outputPool.Subnets[subnet].Touched,
			outputPool.Subnets[subnet].Location, outputPool.Subnets[subnet].Range,
		)
		ch <- prometheus.MustNewConstMetric(
			e.subnetUsed, prometheus.GaugeValue, outputPool.Subnets[subnet].Used,
			outputPool.Subnets[subnet].Location, outputPool.Subnets[subnet].Range,
		)
		ch <- prometheus.MustNewConstMetric(
			e.subnetDefined, prometheus.GaugeValue, outputPool.Subnets[subnet].Defined,
			outputPool.Subnets[subnet].Location, outputPool.Subnets[subnet].Range,
		)
	}

	for sharedNetwork := 0; sharedNetwork < len(outputPool.SharedNetworks); sharedNetwork++ {
		ch <- prometheus.MustNewConstMetric(
			e.sharedNetworkFree, prometheus.GaugeValue, outputPool.SharedNetworks[sharedNetwork].Free,
			outputPool.SharedNetworks[sharedNetwork].Location,
		)
		ch <- prometheus.MustNewConstMetric(
			e.sharedNetworkTouched, prometheus.GaugeValue, outputPool.SharedNetworks[sharedNetwork].Touched,
			outputPool.SharedNetworks[sharedNetwork].Location,
		)
		ch <- prometheus.MustNewConstMetric(
			e.sharedNetworkUsed, prometheus.GaugeValue, outputPool.SharedNetworks[sharedNetwork].Used,
			outputPool.SharedNetworks[sharedNetwork].Location,
		)
		ch <- prometheus.MustNewConstMetric(
			e.sharedNetworkDefined, prometheus.GaugeValue, outputPool.SharedNetworks[sharedNetwork].Defined,
			outputPool.SharedNetworks[sharedNetwork].Location,
		)
	}

	return nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if err := e.collect(ch); err != nil {
		log.Errorf("Error scraping dhcpd-pools: %s", err)
		e.scrapeFailures.Inc()
		e.scrapeFailures.Collect(ch)
	}
	return
}

func init() {
	prometheus.MustRegister(version.NewCollector("iscdhcpd_exporter"))
}

func main() {
	log.AddFlags(kingpin.CommandLine)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	var landingPage = []byte(`<html>
                        <head><title>ISC-DHCPD exporter</title></head>
                        <body>
                        <h1>ISC-DHCPD exporter</h1>
                        <p><a href="` + *metricsPath + `">Metrics</a></p>
                        </body>
                        </html>`)

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	log.Infoln("Starting iscdhcpd_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})
	log.Infoln("Listening on", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
