package installer

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockFileSystem) CopyFile(src, dest string) error {
	args := m.Called(src, dest)
	return args.Error(0)
}

func Test_InstallCore(t *testing.T) {
	t.Parallel()

	mockFS := new(MockFileSystem)

	installer := &WindsurfInstaller{
		templateDir:    "templates",
		localDestDir:   "local",
		localFileName:  ".windsurfrules",
		globalDestDir:  "global",
		globalFileName: "global_rules.md",
		lang:           English,
		fs:             mockFS,
	}

	tests := []struct {
		name string
		args struct {
			installType InstallType
		}
		setup struct {
			mock func()
		}
		want struct {
			err    bool
			errMsg string
		}
	}{
		{
			name: "Successful local installation",
			args: struct{ installType InstallType }{
				installType: Local,
			},
			setup: struct{ mock func() }{
				mock: func() {
					mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(nil)
					mockFS.On("CopyFile", filepath.Join("templates", "local", ".windsurfrules"), filepath.Join("local", ".windsurfrules")).Return(nil)
				},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    false,
				errMsg: "",
			},
		},
		{
			name: "MkdirAll fails",
			args: struct{ installType InstallType }{
				installType: Local,
			},
			setup: struct{ mock func() }{
				mock: func() {
					mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(errors.New("mkdir failed"))
				},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    true,
				errMsg: "failed to create directory",
			},
		},
		{
			name: "CopyFile fails",
			args: struct{ installType InstallType }{
				installType: Local,
			},
			setup: struct{ mock func() }{
				mock: func() {
					mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(nil)
					mockFS.On("CopyFile", filepath.Join("templates", "local", ".windsurfrules"), filepath.Join("local", ".windsurfrules")).Return(errors.New("copy failed"))
				},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    true,
				errMsg: "copy failed",
			},
		},
		{
			name: "Invalid installation type",
			args: struct{ installType InstallType }{
				installType: InstallType(999),
			},
			setup: struct{ mock func() }{
				mock: func() {},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    true,
				errMsg: "unknown installation type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.ExpectedCalls = nil
			tt.setup.mock()

			err := installer.installCore(tt.args.installType)

			if tt.want.err {
				assert.Error(t, err)
				if tt.want.errMsg != "" {
					assert.Contains(t, err.Error(), tt.want.errMsg)
				}
			} else {
				assert.NoError(t, err)
				mockFS.AssertExpectations(t)
			}
		})
	}
}

func Test_Install(t *testing.T) {
	t.Parallel()

	mockFS := new(MockFileSystem)
	installer := &WindsurfInstaller{
		templateDir:    "/template",
		localDestDir:   "/local",
		localFileName:  "local.rules",
		globalDestDir:  "/global",
		globalFileName: "global.rules",
		lang:           English,
		fs:             mockFS,
	}

	tests := []struct {
		name string
		args struct {
			installType InstallType
		}
		setup struct {
			mock func(*MockFileSystem)
		}
		want struct {
			err    bool
			errMsg string
		}
	}{
		{
			name: "Local installation",
			args: struct{ installType InstallType }{
				installType: Local,
			},
			setup: struct{ mock func(*MockFileSystem) }{
				mock: func(mockFS *MockFileSystem) {
					mockFS.On("MkdirAll", "/local", os.FileMode(0755)).Return(nil)
					mockFS.On("CopyFile", filepath.Join("/template", "local", "local.rules"), filepath.Join("/local", "local.rules")).Return(nil)
				},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    false,
				errMsg: "",
			},
		},
		{
			name: "Global installation",
			args: struct{ installType InstallType }{
				installType: Global,
			},
			setup: struct{ mock func(*MockFileSystem) }{
				mock: func(mockFS *MockFileSystem) {
					mockFS.On("MkdirAll", "/global", os.FileMode(0755)).Return(nil)
					mockFS.On("CopyFile", filepath.Join("/template", "global", "global.rules"), filepath.Join("/global", "global.rules")).Return(nil)
				},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    false,
				errMsg: "",
			},
		},
		{
			name: "All installation - success",
			args: struct{ installType InstallType }{
				installType: All,
			},
			setup: struct{ mock func(*MockFileSystem) }{
				mock: func(mockFS *MockFileSystem) {
					mockFS.On("MkdirAll", "/local", os.FileMode(0755)).Return(nil)
					mockFS.On("CopyFile", filepath.Join("/template", "local", "local.rules"), filepath.Join("/local", "local.rules")).Return(nil)
					mockFS.On("MkdirAll", "/global", os.FileMode(0755)).Return(nil)
					mockFS.On("CopyFile", filepath.Join("/template", "global", "global.rules"), filepath.Join("/global", "global.rules")).Return(nil)
				},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    false,
				errMsg: "",
			},
		},
		{
			name: "All installation - local fails",
			args: struct{ installType InstallType }{
				installType: All,
			},
			setup: struct{ mock func(*MockFileSystem) }{
				mock: func(mockFS *MockFileSystem) {
					mockFS.On("MkdirAll", "/local", os.FileMode(0755)).Return(errors.New("mkdir failed"))
				},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    true,
				errMsg: "failed to install local configuration file",
			},
		},
		{
			name: "All installation - global fails",
			args: struct{ installType InstallType }{
				installType: All,
			},
			setup: struct{ mock func(*MockFileSystem) }{
				mock: func(mockFS *MockFileSystem) {
					mockFS.On("MkdirAll", "/local", os.FileMode(0755)).Return(nil)
					mockFS.On("CopyFile", filepath.Join("/template", "local", "local.rules"), filepath.Join("/local", "local.rules")).Return(nil)
					mockFS.On("MkdirAll", "/global", os.FileMode(0755)).Return(nil)
					mockFS.On("CopyFile", filepath.Join("/template", "global", "global.rules"), filepath.Join("/global", "global.rules")).Return(errors.New("copy failed"))
				},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    true,
				errMsg: "failed to install global configuration file",
			},
		},
		{
			name: "Invalid installation type",
			args: struct{ installType InstallType }{
				installType: InstallType(999),
			},
			setup: struct{ mock func(*MockFileSystem) }{
				mock: func(*MockFileSystem) {},
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    true,
				errMsg: "unknown installation type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.ExpectedCalls = nil
			tt.setup.mock(mockFS)

			err := installer.Install(tt.args.installType)

			if tt.want.err {
				assert.Error(t, err)
				if tt.want.errMsg != "" {
					assert.Contains(t, err.Error(), tt.want.errMsg)
				}
			} else {
				assert.NoError(t, err)
				mockFS.AssertExpectations(t)
			}
		})
	}
}

func Test_NewWindsurfInstaller(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args struct {
			templateDir string
			lang        Language
		}
		want struct {
			err    bool
			errMsg string
		}
	}{
		{
			name: "Valid template directory - English",
			args: struct {
				templateDir string
				lang        Language
			}{
				templateDir: "/valid/template/dir",
				lang:        English,
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    false,
				errMsg: "",
			},
		},
		{
			name: "Valid template directory - Japanese",
			args: struct {
				templateDir string
				lang        Language
			}{
				templateDir: "/valid/template/dir",
				lang:        Japanese,
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    false,
				errMsg: "",
			},
		},
		{
			name: "Empty template directory",
			args: struct {
				templateDir string
				lang        Language
			}{
				templateDir: "",
				lang:        English,
			},
			want: struct {
				err    bool
				errMsg string
			}{
				err:    true,
				errMsg: "TEMPLATE_DIR environment variable is not set",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalValue := os.Getenv("TEMPLATE_DIR")
			os.Setenv("TEMPLATE_DIR", tt.args.templateDir)
			defer func() {
				os.Setenv("TEMPLATE_DIR", originalValue)
			}()

			installer, err := NewWindsurfInstaller(tt.args.lang)

			if tt.want.err {
				assert.Error(t, err)
				if tt.want.errMsg != "" {
					assert.Contains(t, err.Error(), tt.want.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, installer)
				assert.Equal(t, tt.args.templateDir, installer.templateDir)
				assert.Equal(t, ".", installer.localDestDir)
				assert.Equal(t, ".windsurfrules", installer.localFileName)
				assert.Contains(t, installer.globalDestDir, ".codeium/windsurf/memories")
				assert.Equal(t, "global_rules.md", installer.globalFileName)
				assert.Equal(t, tt.args.lang, installer.lang)
				assert.IsType(t, &DefaultFileSystem{}, installer.fs)
			}
		})
	}
}
