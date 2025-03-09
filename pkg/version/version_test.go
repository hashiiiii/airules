package version

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup struct {
			version   string
			commit    string
			buildDate string
		}
		want struct {
			containsVersion   bool
			containsCommit    bool
			containsBuildDate bool
		}
	}{
		{
			name: "Get default version information",
			setup: struct {
				version   string
				commit    string
				buildDate string
			}{
				version:   Version,
				commit:    Commit,
				buildDate: BuildDate,
			},
			want: struct {
				containsVersion   bool
				containsCommit    bool
				containsBuildDate bool
			}{
				containsVersion:   true,
				containsCommit:    true,
				containsBuildDate: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origVersion := Version
			origCommit := Commit
			origBuildDate := BuildDate
			defer func() {
				Version = origVersion
				Commit = origCommit
				BuildDate = origBuildDate
			}()

			result := GetVersion()

			assert.NotEmpty(t, result, "GetVersion should not return empty string")
			assert.Equal(t, tt.want.containsVersion, strings.Contains(result, "Version:"), "Should contain Version information")
			assert.Equal(t, tt.want.containsCommit, strings.Contains(result, "Commit:"), "Should contain Commit information")
			assert.Equal(t, tt.want.containsBuildDate, strings.Contains(result, "BuildDate:"), "Should contain BuildDate information")
		})
	}
}
