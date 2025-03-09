package installer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Language_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		lang Language
		want string
	}{
		{
			name: "English language",
			lang: English,
			want: "en",
		},
		{
			name: "Japanese language",
			lang: Japanese,
			want: "ja",
		},
		{
			name: "Unknown language",
			lang: Language(999),
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lang.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetInstallPath(t *testing.T) {
	t.Parallel()

	templateDir := "templates"
	localDestDir := "localDest"
	globalDestDir := "globalDest"
	localFileName := ".windsurfrules"
	globalFileName := "global_rules.md"
	lang := English

	tests := []struct {
		name string
		args struct {
			installType InstallType
		}
		want struct {
			srcPath  string
			destPath string
			destDir  string
			err      bool
		}
	}{
		{
			name: "Local installation type",
			args: struct{ installType InstallType }{
				installType: Local,
			},
			want: struct {
				srcPath  string
				destPath string
				destDir  string
				err      bool
			}{
				srcPath:  "templates/local/.windsurfrules",
				destPath: "localDest/.windsurfrules",
				destDir:  "localDest",
				err:      false,
			},
		},
		{
			name: "Global installation type",
			args: struct{ installType InstallType }{
				installType: Global,
			},
			want: struct {
				srcPath  string
				destPath string
				destDir  string
				err      bool
			}{
				srcPath:  "templates/global/global_rules.md",
				destPath: "globalDest/global_rules.md",
				destDir:  "globalDest",
				err:      false,
			},
		},
		{
			name: "Invalid installation type",
			args: struct{ installType InstallType }{
				installType: InstallType(999),
			},
			want: struct {
				srcPath  string
				destPath string
				destDir  string
				err      bool
			}{
				srcPath:  "",
				destPath: "",
				destDir:  "",
				err:      true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcPath, destPath, destDir, err := GetInstallPath(
				tt.args.installType,
				templateDir,
				localDestDir,
				globalDestDir,
				localFileName,
				globalFileName,
				lang,
			)

			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.srcPath, srcPath)
				assert.Equal(t, tt.want.destPath, destPath)
				assert.Equal(t, tt.want.destDir, destDir)
			}
		})
	}
}
