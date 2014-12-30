package aristalabstatus

import (
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping Ping")
    }
	host := "8.8.8.8"
	if err := PingHost(host); err != nil {
		t.Errorf("Ping to {} failed, should have passed", host)
	}
	host = "8.8.9.9"
	if err := PingHost(host); err == nil {
		t.Errorf("Ping to {} passed, should have failed", host)
	}

}

func TestFetchFlows(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping fetch")
    }
	flows := FetchFlows()
	staticCount := 0
	for _, flow := range flows {
		if strings.HasPrefix(flow.Name, "STATIC") {
			staticCount++
		}
	}
	if staticCount != 4 {
		t.Errorf("Only {}, should be {}", staticCount, 4)
	}
}
