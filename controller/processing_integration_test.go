package controller

import (
	"encoding/json"
	"testing"

	"example.com/gitlab-ci-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// resetMetrics reinitializes the package-level prometheus vectors for test isolation.
func resetMetrics() {
	gitlab_ci_job_number_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_job_number_total",
			Help: "Total number of jobs.",
		},
		[]string{"ci", "project_name", "branch", "user", "status", "job_stage", "job_tags"},
	)
	gitlab_ci_job_time_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_job_time_total",
			Help: "Jobs duration time",
		},
		[]string{"ci", "project_name", "branch", "user", "status", "job_stage", "job_tags"},
	)
	gitlab_ci_pipeline_number_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_pipeline_number_total",
			Help: "Total number of pipelines",
		},
		[]string{"ci", "project_name", "branch", "user", "status"},
	)
	gitlab_ci_pipeline_time_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_pipeline_time_total",
			Help: "Pipeline duration time",
		},
		[]string{"ci", "project_name", "branch", "user", "status"},
	)
	gitlab_ci_pipeline_size = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_pipeline_size",
			Help: "Pipeline Number of jobs",
		},
		[]string{"ci", "project_name", "branch", "user", "status"},
	)
}

func TestInJob_IncrementsCounters(t *testing.T) {
	resetMetrics()
	cfg := &config.Config{
		Default_Branch: "main",
		CI:             "gitlab",
		Branches:       []string{"^main$", "^release-.*$"},
	}
	c := &Controller{Config: cfg}

	job := Job{
		Ref:            "feature/abc",
		ID:             1,
		Stage:          "test",
		Status:         "success",
		Project_path:   "group/project",
		Build_duration: 12,
	}
	job.Runner.Tags = []string{"docker", "small"}
	job.User.Username = "alice"

	b, _ := json.Marshal(job)
	msg, err := c.inJob(b)
	if err != nil || msg != "done" {
		t.Fatalf("expected done with no error, got msg=%q err=%v", msg, err)
	}

	counter, _ := gitlab_ci_job_number_total.GetMetricWithLabelValues(
		"gitlab", "group/project", cfg.Default_Branch, "alice", "success", "test", "docker, small",
	)
	if v := testutil.ToFloat64(counter); v != 1 {
		t.Fatalf("expected job_number_total=1, got %v", v)
	}
	timeCounter, _ := gitlab_ci_job_time_total.GetMetricWithLabelValues(
		"gitlab", "group/project", cfg.Default_Branch, "alice", "success", "test", "docker, small",
	)
	if v := testutil.ToFloat64(timeCounter); v != 12 {
		t.Fatalf("expected job_time_total=12, got %v", v)
	}
}

func TestInPipe_IncrementsCountersAndSize(t *testing.T) {
	resetMetrics()
	cfg := &config.Config{
		Default_Branch: "main",
		CI:             "gitlab",
		Branches:       []string{"^main$", "^release-.*$"},
	}
	c := &Controller{Config: cfg}

	var p Pipeline
	p.Attributes.ID = 100
	p.Attributes.Ref = "release-1.2"
	p.Attributes.Status = "success"
	p.Attributes.Duration = 34
	p.Project.Path = "group/project"
	p.User.Username = "bob"
	p.Jobs = []PJob{
		{Stage: "build", Status: "success", Duration: 10},
		{Stage: "test", Status: "failed", Duration: 5},
	}

	b, _ := json.Marshal(p)
	msg, err := c.inPipe(b)
	if err != nil || msg != "done" {
		t.Fatalf("expected done with no error, got msg=%q err=%v", msg, err)
	}

	pipeCounter, _ := gitlab_ci_pipeline_number_total.GetMetricWithLabelValues(
		"gitlab", "group/project", "release-1.2", "bob", "success",
	)
	if v := testutil.ToFloat64(pipeCounter); v != 1 {
		t.Fatalf("expected pipeline_number_total=1, got %v", v)
	}
	pipeTime, _ := gitlab_ci_pipeline_time_total.GetMetricWithLabelValues(
		"gitlab", "group/project", "release-1.2", "bob", "success",
	)
	if v := testutil.ToFloat64(pipeTime); v != 34 {
		t.Fatalf("expected pipeline_time_total=34, got %v", v)
	}

	// size: one success and one failed
	sizeSuccess, _ := gitlab_ci_pipeline_size.GetMetricWithLabelValues(
		"gitlab", "group/project", "release-1.2", "bob", "success",
	)
	sizeFailed, _ := gitlab_ci_pipeline_size.GetMetricWithLabelValues(
		"gitlab", "group/project", "release-1.2", "bob", "failed",
	)
	if v := testutil.ToFloat64(sizeSuccess); v != 1 {
		t.Fatalf("expected pipeline_size[success]=1, got %v", v)
	}
	if v := testutil.ToFloat64(sizeFailed); v != 1 {
		t.Fatalf("expected pipeline_size[failed]=1, got %v", v)
	}
}
