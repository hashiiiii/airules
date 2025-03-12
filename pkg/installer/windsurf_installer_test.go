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

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	args := m.Called(path, data, perm)
	return args.Error(0)
}

func TestWindsurfInstaller_InstallLocal(t *testing.T) {
	t.Parallel()

	mockFS := new(MockFileSystem)
	installer := &WindsurfInstaller{
		localDestDir:   "local",
		localFileName:  ".windsurfrules",
		globalDestDir:  "global",
		globalFileName: "global_rules.md",
		fs:             mockFS,
	}

	// テストデータ
	rulePaths := []string{"path/to/rule1.md", "path/to/rule2.md"}

	// 成功ケース
	t.Run("Success", func(t *testing.T) {
		mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(nil)
		mockFS.On("ReadFile", "path/to/rule1.md").Return([]byte("rule1 content"), nil)
		mockFS.On("ReadFile", "path/to/rule2.md").Return([]byte("rule2 content"), nil)
		mockFS.On("WriteFile", filepath.Join("local", ".windsurfrules"), mock.Anything, os.FileMode(0644)).Return(nil)

		err := installer.installLocal(rulePaths)
		assert.NoError(t, err)
		mockFS.AssertExpectations(t)
	})

	// MkdirAll失敗ケース
	t.Run("MkdirAll fails", func(t *testing.T) {
		mockFS.ExpectedCalls = nil
		mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(errors.New("mkdir failed"))

		err := installer.installLocal(rulePaths)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create directory")
	})

	// ReadFile失敗ケース
	t.Run("ReadFile fails", func(t *testing.T) {
		mockFS.ExpectedCalls = nil
		mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(nil)
		mockFS.On("ReadFile", "path/to/rule1.md").Return(nil, errors.New("read failed"))

		err := installer.installLocal(rulePaths)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read rule file")
	})

	// WriteFile失敗ケース
	t.Run("WriteFile fails", func(t *testing.T) {
		mockFS.ExpectedCalls = nil
		mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(nil)
		mockFS.On("ReadFile", "path/to/rule1.md").Return([]byte("rule1 content"), nil)
		mockFS.On("ReadFile", "path/to/rule2.md").Return([]byte("rule2 content"), nil)
		mockFS.On("WriteFile", filepath.Join("local", ".windsurfrules"), mock.Anything, os.FileMode(0644)).Return(errors.New("write failed"))

		err := installer.installLocal(rulePaths)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write to")
	})
}

func Test_NewWindsurfInstaller(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want struct {
			err    bool
			errMsg string
		}
	}{
		{
			name: "Valid installer creation",
			want: struct {
				err    bool
				errMsg string
			}{
				err:    false,
				errMsg: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installer, err := NewWindsurfInstaller()

			if tt.want.err {
				assert.Error(t, err)
				if tt.want.errMsg != "" {
					assert.Contains(t, err.Error(), tt.want.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, installer)
				assert.Equal(t, ".", installer.localDestDir)
				assert.Equal(t, ".windsurfrules", installer.localFileName)
				assert.Contains(t, installer.globalDestDir, filepath.Join(".codeium", "windsurf", "memories"))
				assert.Equal(t, "global_rules.md", installer.globalFileName)
				assert.IsType(t, &DefaultFileSystem{}, installer.fs)
			}
		})
	}
}
