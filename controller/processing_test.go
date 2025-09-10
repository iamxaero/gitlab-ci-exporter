package controller

import (
	"testing"

	"example.com/gitlab-ci-exporter/config"
)

func TestGroupingBranch_DefaultWhenNoMatch(t *testing.T) {
	cfg := &config.Config{
		Default_Branch: "main",
		CI:             "gitlab",
		Branches:       []string{"release/.*", "hotfix/.*"},
	}
	c := &Controller{Config: cfg}

	got := c.GroupingBranch("feature/awesome")
	if got != "main" {
		t.Fatalf("expected default branch 'main', got %q", got)
	}
}

func TestGroupingBranch_MatchReturnsOriginalRef(t *testing.T) {
	cfg := &config.Config{
		Default_Branch: "main",
		CI:             "gitlab",
		Branches:       []string{"feature/.*", "release/.*"},
	}
	c := &Controller{Config: cfg}

	got := c.GroupingBranch("feature/awesome")
	if got != "feature/awesome" {
		t.Fatalf("expected original ref when matched, got %q", got)
	}
}

func TestGroupingBranch_FirstMatchStops(t *testing.T) {
	cfg := &config.Config{
		Default_Branch: "develop",
		CI:             "gitlab",
		Branches:       []string{"feat-.*", "feat-123.*"},
	}
	c := &Controller{Config: cfg}

	got := c.GroupingBranch("feat-123-add-tests")
	if got != "feat-123-add-tests" {
		t.Fatalf("expected matched ref, got %q", got)
	}
}
