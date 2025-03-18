package installer

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCursorInstaller_InstallLocal(t *testing.T) {
	t.Parallel()

	mockFS := new(MockFileSystem)
	installer := &CursorInstaller{
		BaseInstaller: BaseInstaller{
			localDestDir:   "local",
			localFileName:  "project_rules.mdc",
			globalDestDir:  "global",
			globalFileName: "global_rules.mdc",
			editor:         "cursor",
			fs:             mockFS,
		},
	}

	// テストデータ
	testContent := []byte("test content")

	t.Run("Success", func(t *testing.T) {
		mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(nil).Once()
		mockFS.On("ReadFile", "path1").Return(testContent, nil).Once()
		mockFS.On("WriteFile", filepath.Join("local", "project_rules.mdc"), mock.Anything, os.FileMode(0644)).Return(nil).Once()

		err := installer.BaseInstaller.installLocal([]string{"path1"})
		assert.NoError(t, err)
		mockFS.AssertExpectations(t)
	})

	t.Run("MkdirAll_fails", func(t *testing.T) {
		mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(errors.New("mkdir error")).Once()

		err := installer.BaseInstaller.installLocal([]string{"path1"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create directory")
		mockFS.AssertExpectations(t)
	})

	t.Run("ReadFile_fails", func(t *testing.T) {
		mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(nil).Once()
		mockFS.On("ReadFile", "path1").Return(nil, errors.New("read error")).Once()

		err := installer.BaseInstaller.installLocal([]string{"path1"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read rule file")
		mockFS.AssertExpectations(t)
	})

	t.Run("WriteFile_fails", func(t *testing.T) {
		mockFS.On("MkdirAll", "local", os.FileMode(0755)).Return(nil).Once()
		mockFS.On("ReadFile", "path1").Return(testContent, nil).Once()
		mockFS.On("WriteFile", filepath.Join("local", "project_rules.mdc"), mock.Anything, os.FileMode(0644)).Return(errors.New("write error")).Once()

		err := installer.BaseInstaller.installLocal([]string{"path1"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write to")
		mockFS.AssertExpectations(t)
	})
}

func Test_NewCursorInstaller(t *testing.T) {
	t.Parallel()

	t.Run("Valid_installer_creation", func(t *testing.T) {
		installer, err := NewCursorInstaller()
		assert.NoError(t, err)
		assert.NotNil(t, installer)
		assert.NotEmpty(t, installer.BaseInstaller.localDestDir)
		assert.NotEmpty(t, installer.BaseInstaller.globalDestDir)
		assert.Equal(t, "project_rules.mdc", installer.BaseInstaller.localFileName)
		assert.Equal(t, "global_rules.mdc", installer.BaseInstaller.globalFileName)
		assert.Equal(t, "cursor", installer.BaseInstaller.editor)
		assert.NotNil(t, installer.BaseInstaller.fs)
	})
}
