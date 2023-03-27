package source

import (
	"teredix/pkg/config"
	"teredix/pkg/source/scanner"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSources(t *testing.T) {
	appConfig := &config.AppConfig{
		Sources: map[string]config.Source{
			"source1": {
				Type: "file_system",
				Configuration: map[string]string{
					"root_directory": "/path/to/directory",
				},
			},
		},
	}

	sources := BuildSources(appConfig)

	fsScanner := scanner.NewFsScanner("fs-scanner_1", "/path/to/directory", map[string]string{})

	expectedSources := []Source{
		{
			Name:    "source1",
			Scanner: &fsScanner,
		},
	}

	assert.Equal(t, expectedSources, sources)
}