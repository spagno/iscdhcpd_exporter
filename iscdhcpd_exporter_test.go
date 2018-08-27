package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	dhcpdPools = `{
   "subnets": [{
      "location":"All networks", "range":"192.168.1.10 - 192.168.1.100", "defined":91, "used":1, "touched":55, "free":90 }, { "location":"All networks", "range":"192.168.2.10 - 192.168.2.100", "defined":91, "used":5, "touched":25, "free":86 }, { "location":"All networks", "range":"192.168.3.10 - 192.168.3.100", "defined":91, "used":0, "touched":0, "free":91 }], "shared-networks": [], "summary": {"location":"All networks", "defined":1, "used":2, "touched":3, "free":4}}`
)

func checkPoolsStatus(t *testing.T, status string, metricCount int) {
	e := NewExporter()
	ch := make(chan prometheus.Metric)
	var f ToTest
	f.Output = []byte(dhcpdPools)

	go func() {
		defer close(ch)
		e.Collect(ch)
	}()

	for i := 1; i <= metricCount; i++ {
		m := <-ch
		if m == nil {
			t.Error("expected metric but got nil")
		}
	}
	extraMetrics := 0
	for <-ch != nil {
		extraMetrics++
	}
	if extraMetrics > 0 {
		t.Errorf("expected closed channel, got %d extra metrics", extraMetrics)
	}
}

func TestPoolsStatus(t *testing.T) {
	checkPoolsStatus(t, dhcpdPools, 16)
}
