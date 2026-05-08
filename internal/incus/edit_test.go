package incus

import (
	"strings"
	"testing"

	"github.com/lxc/incus/v7/shared/api"
)

func TestInstanceConfigEditContentHidesExpandedFields(t *testing.T) {
	content, err := instanceConfigEditContent(&api.Instance{
		Name:            "ubuntu",
		ExpandedConfig:  map[string]string{"limits.cpu": "2"},
		ExpandedDevices: map[string]map[string]string{"root": {"path": "/"}},
		InstancePut: api.InstancePut{
			Config: map[string]string{"boot.autostart": "true"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected edit content error: %v", err)
	}

	text := string(content)
	for _, want := range []string{"YAML representation", "name: ubuntu", "boot.autostart: \"true\""} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected edit content to contain %q\ncontent:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{"expanded_config", "expanded_devices", "limits.cpu"} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("expected edit content not to contain %q\ncontent:\n%s", unwanted, text)
		}
	}
}

func TestParseInstanceConfigEditReadsWritableFields(t *testing.T) {
	updated, err := parseInstanceConfigEdit([]byte(`config:
  boot.autostart: "true"
description: edited instance
profiles:
  - default
`))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if updated.Description != "edited instance" {
		t.Fatalf("unexpected description: %q", updated.Description)
	}
	if updated.Config["boot.autostart"] != "true" {
		t.Fatalf("unexpected config: %v", updated.Config)
	}
	if len(updated.Profiles) != 1 || updated.Profiles[0] != "default" {
		t.Fatalf("unexpected profiles: %v", updated.Profiles)
	}
}

func TestParseInstanceConfigEditRejectsInvalidYAML(t *testing.T) {
	if _, err := parseInstanceConfigEdit([]byte("config: [")); err == nil {
		t.Fatal("expected invalid YAML error")
	}
}
