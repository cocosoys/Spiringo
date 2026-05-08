package metrics

import (
	"strings"
	"testing"
)

// 中文：TestRegistryWritesPrometheusCountersAndSummaries 验证相关行为符合预期。
// English: TestRegistryWritesPrometheusCountersAndSummaries verifies the related behavior.
func TestRegistryWritesPrometheusCountersAndSummaries(t *testing.T) {
	registry := NewRegistry("app")
	labels := Labels{"method": "GET", "path": "/health", "status": "200"}
	registry.IncCounter("http_requests_total", labels)
	registry.AddCounter("http_requests_total", labels, 2)
	registry.ObserveSummary("http_request_duration_seconds", labels, 0.1)
	registry.ObserveSummary("http_request_duration_seconds", labels, 0.2)

	var out strings.Builder
	if err := registry.WritePrometheus(&out); err != nil {
		t.Fatalf("WritePrometheus returned error: %v", err)
	}
	text := out.String()
	assertContains(t, text, "# TYPE app_http_requests_total counter")
	assertContains(t, text, `app_http_requests_total{method="GET",path="/health",status="200"} 3`)
	assertContains(t, text, "# TYPE app_http_request_duration_seconds summary")
	assertContains(t, text, `app_http_request_duration_seconds_sum{method="GET",path="/health",status="200"} 0.3`)
	assertContains(t, text, `app_http_request_duration_seconds_count{method="GET",path="/health",status="200"} 2`)
}

// 中文：TestRegistryEscapesLabels 验证相关行为符合预期。
// English: TestRegistryEscapesLabels verifies the related behavior.
func TestRegistryEscapesLabels(t *testing.T) {
	registry := NewRegistry("")
	registry.IncCounter("events_total", Labels{"bad-label": "line\n\"quoted\""})

	var out strings.Builder
	if err := registry.WritePrometheus(&out); err != nil {
		t.Fatalf("WritePrometheus returned error: %v", err)
	}
	assertContains(t, out.String(), `bad_label="line\n\"quoted\""`)
}

// 中文：TestRegistrySnapshotAndReport 验证相关行为符合预期。
// English: TestRegistrySnapshotAndReport verifies the related behavior.
func TestRegistrySnapshotAndReport(t *testing.T) {
	registry := NewRegistry("spiringo")
	registry.AddCounter("jobs_total", Labels{"status": "ok"}, 3)
	registry.ObserveSummary("job_duration_seconds", Labels{"kind": "email"}, 0.25)
	registry.ObserveSummary("job_duration_seconds", Labels{"kind": "email"}, 0.75)

	snapshot := registry.Snapshot()
	if len(snapshot.Counters) != 1 || snapshot.Counters[0].Name != "spiringo_jobs_total" {
		t.Fatalf("unexpected counters: %+v", snapshot.Counters)
	}
	if len(snapshot.Summaries) != 1 || snapshot.Summaries[0].Average != 0.5 {
		t.Fatalf("unexpected summaries: %+v", snapshot.Summaries)
	}

	var out strings.Builder
	if err := registry.WriteReport(&out); err != nil {
		t.Fatalf("WriteReport returned error: %v", err)
	}
	text := out.String()
	assertContains(t, text, "# Metrics Report")
	assertContains(t, text, "`spiringo_jobs_total`")
	assertContains(t, text, "`status=ok`")
	assertContains(t, text, "`spiringo_job_duration_seconds`")
	assertContains(t, text, "0.500000")
}

// 中文：assertContains 执行当前包中的对应流程。
// English: assertContains executes the corresponding workflow in this package.
func assertContains(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		t.Fatalf("output missing %q:\n%s", want, text)
	}
}
