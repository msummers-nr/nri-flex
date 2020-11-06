// +build integration
// +build windows

package integration_test

import (
	"path/filepath"
	"testing"

	sdk "github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/runtime"
	"github.com/stretchr/testify/require"
)

var configDirPath = filepath.Join("configs", "windows")

func Test_WindowsCommands_ReturnsData(t *testing.T) {
	configDirPath := filepath.Join("configs", "windows")

	load.Refresh()

	i, _ := sdk.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// set file to load
	load.Args.ConfigFile = filepath.Join(configDirPath, "windows-cmd-test.yml")

	// when
	r := runtime.GetDefaultRuntime()
	err := runtime.RunFlex(r)
	require.NoError(t, err)

	//then
	metricsSet := load.Entity.Metrics
	require.NotEmpty(t, metricsSet)

	for _, ms := range metricsSet {
		if ms.Metrics["event_type"] == "flexStatusSample" {
			continue
		}
		require.NotNil(t, ms.Metrics["status"], "status")
		require.NotNil(t, ms.Metrics["name"], "name")
		require.NotNil(t, ms.Metrics["displayname"], "displayname")
	}

	// check for a specific service, because Flex ingests everything, even output "garbage"
	// any Windows version should always have the Themes service, so check for that
	var found bool
	for _, ms := range metricsSet {
		if ms.Metrics["name"] == "Themes" {
			found = true
		}
	}

	require.Truef(t, found, "didn't find the 'Themes' service. check that the configuration is correct")
}

func Test_WindowsCommands_ConfigsFolder(t *testing.T) {
	load.Refresh()

	i, _ := sdk.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// set configurations folder
	load.Args.ConfigDir = configDirPath

	// when
	r := runtime.GetDefaultRuntime()
	err := runtime.RunFlex(r)
	require.NoError(t, err)

	//then
	metricsSet := load.Entity.Metrics
	require.NotEmpty(t, metricsSet)

	for _, ms := range metricsSet {
		// ignore all other samples from other configs in the folder
		if ms.Metrics["event_type"] != "windowsServiceListSample" {
			continue
		}
		require.NotNil(t, ms.Metrics["status"], "status")
		require.NotNil(t, ms.Metrics["name"], "name")
		require.NotNil(t, ms.Metrics["displayname"], "displayname")
	}

	// check for a specific service, because Flex ingests everything, even output "garbage"
	// any Windows version should always have the Themes service, so check for that
	var found bool
	for _, ms := range metricsSet {
		if ms.Metrics["name"] == "Themes" {
			found = true
		}
	}

	require.Truef(t, found, "didn't find the 'Themes' service. check that the configuration is correct")
}
