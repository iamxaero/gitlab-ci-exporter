package controller

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	gitlab_ci_job_number_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_job_number_total",
			Help: "Total number of jobs.",
		},
		[]string{
			"ci",
			"project_name",
			"branch",
			"user",
			"status",
			"job_stage",
			"job_tags",
		},
	)
	gitlab_ci_job_time_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_job_time_total",
			Help: "Jobs duration time",
		},
		[]string{
			"ci",
			"project_name",
			"branch",
			"user",
			"status",
			"job_stage",
			"job_tags",
		},
	)
	gitlab_ci_pipeline_number_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_pipeline_number_total",
			Help: "Total number of pipelines",
		},
		[]string{
			"ci",
			"project_name",
			"branch",
			"user",
			"status",
		},
	)
	gitlab_ci_pipeline_time_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_pipeline_time_total",
			Help: "Pipeline duration time",
		},
		[]string{
			"ci",
			"project_name",
			"branch",
			"user",
			"status",
		},
	)
	gitlab_ci_pipeline_size = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_ci_pipeline_size",
			Help: "Pipeline Number of jobs",
		},
		[]string{
			"ci",
			"project_name",
			"branch",
			"user",
			"status",
		},
	)
)

func (c *Controller) PromRegister() {
	// register metrics
	prometheus.MustRegister(gitlab_ci_job_number_total)
	prometheus.MustRegister(gitlab_ci_job_time_total)
	prometheus.MustRegister(gitlab_ci_pipeline_number_total)
	prometheus.MustRegister(gitlab_ci_pipeline_time_total)
	prometheus.MustRegister(gitlab_ci_pipeline_size)
}

func (c *Controller) GroupingBranch(ref string) string {
	set_ref := c.Config.Default_Branch
BranchesLoop:
	for _, v := range c.Config.Branches {
		matched, _ := regexp.MatchString(v, ref)
		switch matched {
		case false:
			continue
		case true:
			set_ref = ref
			break BranchesLoop
		}
	}
	ref = set_ref
	return ref
}

// Pipelines processing
func (c *Controller) inPipe(data []byte) (string, error) {
	msg := "skipped"
	var err error
	var pipe Pipeline

	err = json.Unmarshal(data, &pipe)
	if err != nil {
		return "failed to decode json", err
	}

	set_ref := c.GroupingBranch(pipe.Attributes.Ref)
	pipe.Attributes.Ref = set_ref

	// Set Pipeline values
	gitlab_ci_pipeline_number_total.WithLabelValues(
		c.Config.CI,
		pipe.Project.Path,
		pipe.Attributes.Ref,
		pipe.User.Username,
		pipe.Attributes.Status,
	).Add(1)
	gitlab_ci_pipeline_time_total.WithLabelValues(
		c.Config.CI,
		pipe.Project.Path,
		pipe.Attributes.Ref,
		pipe.User.Username,
		pipe.Attributes.Status,
	).Add(float64(pipe.Attributes.Duration))
	// Count jobs from pipeline data
	for _, v := range pipe.Jobs {
		gitlab_ci_pipeline_size.WithLabelValues(
			c.Config.CI,
			pipe.Project.Path,
			pipe.Attributes.Ref,
			pipe.User.Username,
			v.Status,
		).Add(1)
	}
	log.Debugf("pipeline is processed ID: %v Status: %v", pipe.Attributes.ID, pipe.Attributes.Status)
	msg = "done"

	return msg, nil
}

// Jobs processing
func (c *Controller) inJob(data []byte) (string, error) {
	msg := "skipped"
	var err error
	var job Job

	err = json.Unmarshal(data, &job)
	if err != nil {
		return "failed to decode json", err
	}

	set_ref := c.GroupingBranch(job.Ref)
	job.Ref = set_ref

	// Set Job values
	job_runner := strings.Join(job.Runner.Tags, ", ")
	gitlab_ci_job_number_total.WithLabelValues(
		c.Config.CI,
		strings.ReplaceAll(job.Project_path, " ", ""),
		job.Ref,
		job.User.Username,
		job.Status,
		job.Stage,
		job_runner,
	).Add(1)
	gitlab_ci_job_time_total.WithLabelValues(
		c.Config.CI,
		strings.ReplaceAll(job.Project_path, " ", ""),
		job.Ref,
		job.User.Username,
		job.Status,
		job.Stage,
		job_runner,
	).Add(float64(job.Build_duration))

	log.Debugf("job is processed ID: %v Status: %v", job.ID, job.Status)
	msg = "done"
	return msg, nil
}
