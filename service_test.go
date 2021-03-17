package metrics

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type results struct {
	returnValue []versionInformation
	returnError error
}

func (s *results) get() ([]versionInformation, error) {
	return s.returnValue, s.returnError
}

func newTestService(metricName string) (service, *results) {
	res := new(results)
	return NewService(
		metricName,
		"This is not a helpful description.",
		res.get,
	), res
}

func formatVersion(metricName string, v versionInformation) string {
	return fmt.Sprintf("%s{id=%q,min_version=%q,version=%q} 1", metricName, v.ID, v.MinVersion, v.Version)
}

func getMetrics(t *testing.T, u, metricName string) []string {
	res, err := http.Get(u)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var lines []string
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		if line := scanner.Text(); strings.HasPrefix(line, metricName) {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	return lines
}

func TestRefresh(t *testing.T) {
	const metricName = "my_example_metric"

	s, expectedVersions := newTestService(metricName)

	r := prometheus.NewRegistry()
	r.MustRegister(s)

	// http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	ts := httptest.NewServer(promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	defer ts.Close()

	expectedVersions.returnValue = []versionInformation{
		{ID: "1", MinVersion: "1.1", Version: "1.99"},
	}
	if err := s.Refresh(); err != nil {
		t.Fatal(err)
	}

	{
		metrics := getMetrics(t, ts.URL, metricName)
		if want, have := len(metrics), len(expectedVersions.returnValue); want != have {
			t.Errorf("expected %d metrics, got %d", want, have)
		}

		for _, rv := range expectedVersions.returnValue {
			var found bool
			want := formatVersion(metricName, rv)
			for _, have := range metrics {
				if have == want {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected '%s', not found", want)
			}
		}
	}

	expectedVersions.returnValue = []versionInformation{
		{ID: "1", MinVersion: "1.1", Version: "1.202"},
		{ID: "2", MinVersion: "2.1", Version: "2.99"},
	}
	if err := s.Refresh(); err != nil {
		t.Fatal(err)
	}

	{
		metrics := getMetrics(t, ts.URL, metricName)
		if want, have := len(metrics), len(expectedVersions.returnValue); want != have {
			t.Errorf("expected %d metrics, got %d", want, have)
		}

		for _, rv := range expectedVersions.returnValue {
			var found bool
			want := formatVersion(metricName, rv)
			for _, have := range metrics {
				if have == want {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected '%s', not found", want)
			}
		}
	}

	expectedVersions.returnValue = []versionInformation{
		{ID: "2", MinVersion: "2.3", Version: "2.101"},
	}
	if err := s.Refresh(); err != nil {
		t.Fatal(err)
	}

	{
		metrics := getMetrics(t, ts.URL, metricName)
		if want, have := len(metrics), len(expectedVersions.returnValue); want != have {
			t.Errorf("expected %d metrics, got %d", want, have)
		}

		for _, rv := range expectedVersions.returnValue {
			var found bool
			want := formatVersion(metricName, rv)
			for _, have := range metrics {
				if have == want {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected '%s', not found", want)
			}
		}
	}
}
