package installer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetInstallPath(t *testing.T) {
	t.Parallel()

	templateDir := "templates"
	localDestDir := "localDest"
	globalDestDir := "globalDest"
	localFileName := ".windsurfrules"
	globalFileName := "global_rules.md"

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
			name: "Local installation",
			args: struct {
				installType InstallType
			}{
				installType: Local,
			},
			want: struct {
				srcPath  string
				destPath string
				destDir  string
				err      bool
			}{
				srcPath:  filepath.Join(templateDir, "local", localFileName),
				destPath: filepath.Join(localDestDir, localFileName),
				destDir:  localDestDir,
				err:      false,
			},
		},
		{
			name: "Global installation",
			args: struct {
				installType InstallType
			}{
				installType: Global,
			},
			want: struct {
				srcPath  string
				destPath string
				destDir  string
				err      bool
			}{
				srcPath:  filepath.Join(templateDir, "global", globalFileName),
				destPath: filepath.Join(globalDestDir, globalFileName),
				destDir:  globalDestDir,
				err:      false,
			},
		},
		{
			name: "Unknown installation type",
			args: struct {
				installType InstallType
			}{
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
			srcPath, destPath, destDir, err := GetInstallPath(tt.args.installType, templateDir, localDestDir, globalDestDir, localFileName, globalFileName)

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
