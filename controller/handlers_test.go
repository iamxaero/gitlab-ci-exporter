package controller

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/gitlab-ci-exporter/config"
)

func newTestController() *Controller {
	cfg := &config.Config{
		Default_Branch: "main",
		CI:             "gitlab",
		Branches:       []string{"^main$", "^release-.*$"},
	}
	return &Controller{Config: cfg}
}

func TestHealth(t *testing.T) {
	c := newTestController()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	c.Health(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	b, _ := io.ReadAll(res.Body)
	if !bytes.Contains(b, []byte("Gitlab CI Exporter is OK")) {
		t.Fatalf("unexpected body: %q", string(b))
	}
}

func TestWebhook_NoHeaderSkips(t *testing.T) {
	c := newTestController()
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString("{}"))
	rec := httptest.NewRecorder()

	c.Webhook(rec, req)

	if body := rec.Body.String(); body != "skipped" {
		t.Fatalf("expected 'skipped', got %q", body)
	}
}

func TestWebhook_JobHook(t *testing.T) {
	c := newTestController()
	jobJSON := `{"ref":"main","build_id":1,"build_stage":"test","build_status":"success","project_name":"group/project","build_duration":2,"runner":{"tags":["x"]},"user":{"username":"bob"}}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(jobJSON))
	req.Header.Set("X-Gitlab-Event", "Job Hook")
	rec := httptest.NewRecorder()

	c.Webhook(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if body := rec.Body.String(); body != "done" {
		t.Fatalf("expected 'done', got %q", body)
	}
}

func TestWebhook_PipelineHook(t *testing.T) {
	c := newTestController()
	pipeJSON := `{"object_attributes":{"id":10,"ref":"main","status":"success","duration":3},"project":{"path_with_namespace":"group/project"},"user":{"username":"alice"},"builds":[]}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(pipeJSON))
	req.Header.Set("X-Gitlab-Event", "Pipeline Hook")
	rec := httptest.NewRecorder()

	c.Webhook(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if body := rec.Body.String(); body != "done" {
		t.Fatalf("expected 'done', got %q", body)
	}
}
