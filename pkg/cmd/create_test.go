/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"helm.sh/helm/v4/internal/test/ensure"
	chart "helm.sh/helm/v4/pkg/chart/v2"
	"helm.sh/helm/v4/pkg/chart/v2/loader"
	chartutil "helm.sh/helm/v4/pkg/chart/v2/util"
	"helm.sh/helm/v4/pkg/helmpath"
)

func TestCreateCmd(t *testing.T) {
	t.Chdir(t.TempDir())
	ensure.HelmHome(t)
	cname := "testchart"

	// Run a create
	if _, _, err := executeActionCommand("create " + cname); err != nil {
		t.Fatalf("Failed to run create: %s", err)
	}

	// Test that the chart is there
	if fi, err := os.Stat(cname); err != nil {
		t.Fatalf("no chart directory: %s", err)
	} else if !fi.IsDir() {
		t.Fatalf("chart is not directory")
	}

	c, err := loader.LoadDir(cname)
	if err != nil {
		t.Fatal(err)
	}

	if c.Name() != cname {
		t.Errorf("Expected %q name, got %q", cname, c.Name())
	}
	if c.Metadata.APIVersion != chart.APIVersionV2 {
		t.Errorf("Wrong API version: %q", c.Metadata.APIVersion)
	}
}

func TestCreateStarterCmd(t *testing.T) {
	t.Chdir(t.TempDir())
	ensure.HelmHome(t)
	cname := "testchart"
	defer resetEnv()()
	// Create a starter.
	starterchart := helmpath.DataPath("starters")
	os.MkdirAll(starterchart, 0o755)
	if dest, err := chartutil.Create("starterchart", starterchart); err != nil {
		t.Fatalf("Could not create chart: %s", err)
	} else {
		t.Logf("Created %s", dest)
	}
	tplpath := filepath.Join(starterchart, "starterchart", "templates", "foo.tpl")
	if err := os.WriteFile(tplpath, []byte("test"), 0o644); err != nil {
		t.Fatalf("Could not write template: %s", err)
	}

	// Run a create
	if _, _, err := executeActionCommand(fmt.Sprintf("create --starter=starterchart %s", cname)); err != nil {
		t.Errorf("Failed to run create: %s", err)
		return
	}

	// Test that the chart is there
	if fi, err := os.Stat(cname); err != nil {
		t.Fatalf("no chart directory: %s", err)
	} else if !fi.IsDir() {
		t.Fatalf("chart is not directory")
	}

	c, err := loader.LoadDir(cname)
	if err != nil {
		t.Fatal(err)
	}

	if c.Name() != cname {
		t.Errorf("Expected %q name, got %q", cname, c.Name())
	}
	if c.Metadata.APIVersion != chart.APIVersionV2 {
		t.Errorf("Wrong API version: %q", c.Metadata.APIVersion)
	}

	expectedNumberOfTemplates := 10
	if l := len(c.Templates); l != expectedNumberOfTemplates {
		t.Errorf("Expected %d templates, got %d", expectedNumberOfTemplates, l)
	}

	found := false
	for _, tpl := range c.Templates {
		if tpl.Name == "templates/foo.tpl" {
			found = true
			if data := string(tpl.Data); data != "test" {
				t.Errorf("Expected template 'test', got %q", data)
			}
		}
	}
	if !found {
		t.Error("Did not find foo.tpl")
	}
}

func TestCreateStarterAbsoluteCmd(t *testing.T) {
	t.Chdir(t.TempDir())
	defer resetEnv()()
	ensure.HelmHome(t)
	cname := "testchart"

	// Create a starter.
	starterchart := helmpath.DataPath("starters")
	os.MkdirAll(starterchart, 0o755)
	if dest, err := chartutil.Create("starterchart", starterchart); err != nil {
		t.Fatalf("Could not create chart: %s", err)
	} else {
		t.Logf("Created %s", dest)
	}
	tplpath := filepath.Join(starterchart, "starterchart", "templates", "foo.tpl")
	if err := os.WriteFile(tplpath, []byte("test"), 0o644); err != nil {
		t.Fatalf("Could not write template: %s", err)
	}

	starterChartPath := filepath.Join(starterchart, "starterchart")

	// Run a create
	if _, _, err := executeActionCommand(fmt.Sprintf("create --starter=%s %s", starterChartPath, cname)); err != nil {
		t.Errorf("Failed to run create: %s", err)
		return
	}

	// Test that the chart is there
	if fi, err := os.Stat(cname); err != nil {
		t.Fatalf("no chart directory: %s", err)
	} else if !fi.IsDir() {
		t.Fatalf("chart is not directory")
	}

	c, err := loader.LoadDir(cname)
	if err != nil {
		t.Fatal(err)
	}

	if c.Name() != cname {
		t.Errorf("Expected %q name, got %q", cname, c.Name())
	}
	if c.Metadata.APIVersion != chart.APIVersionV2 {
		t.Errorf("Wrong API version: %q", c.Metadata.APIVersion)
	}

	expectedNumberOfTemplates := 10
	if l := len(c.Templates); l != expectedNumberOfTemplates {
		t.Errorf("Expected %d templates, got %d", expectedNumberOfTemplates, l)
	}

	found := false
	for _, tpl := range c.Templates {
		if tpl.Name == "templates/foo.tpl" {
			found = true
			if data := string(tpl.Data); data != "test" {
				t.Errorf("Expected template 'test', got %q", data)
			}
		}
	}
	if !found {
		t.Error("Did not find foo.tpl")
	}
}

func TestCreateFileCompletion(t *testing.T) {
	checkFileCompletion(t, "create", true)
	checkFileCompletion(t, "create myname", false)
}
