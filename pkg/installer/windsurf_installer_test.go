package installer

import (
	"os"
	"path/filepath"
	"testing"
)

// インストーラーのテスト用ヘルパー関数
func setupInstallerTest(t *testing.T) (string, string, string) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	templateDir := filepath.Join(tmpDir, "templates")
	destDir := filepath.Join(tmpDir, "dest")
	
	// テスト用のディレクトリを作成
	err := os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("テンプレートディレクトリの作成に失敗: %v", err)
	}
	
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("インストール先ディレクトリの作成に失敗: %v", err)
	}
	
	return tmpDir, templateDir, destDir
}

// Windsurfインストーラーのローカルインストールテスト
func TestWindsurfInstaller_InstallLocal(t *testing.T) {
	// テスト環境のセットアップ
	_, templateDir, destDir := setupInstallerTest(t)
	
	// テンプレートファイルを作成
	templateFilePath := filepath.Join(templateDir, "cascade.local.json")
	templateContent := []byte(`{"test": "windsurf-config"}`)
	err := os.WriteFile(templateFilePath, templateContent, 0644)
	if err != nil {
		t.Fatalf("テンプレートファイルの作成に失敗: %v", err)
	}
	
	// インストーラーを作成
	installer := NewWindsurfInstaller()
	
	// テスト用にインストール先ディレクトリとテンプレートディレクトリを上書き
	installer.templateDir = templateDir
	installer.localDestDir = destDir
	
	// ローカルファイルのインストールをテスト
	err = installer.InstallLocal()
	if err != nil {
		t.Fatalf("ローカルファイルのインストールに失敗: %v", err)
	}
	
	// インストールされたファイルが存在するか確認
	destFile := filepath.Join(destDir, installer.localFileName)
	exists := fileExists(destFile)
	if !exists {
		t.Errorf("ファイルがインストールされていません: %s", destFile)
	}
	
	// ファイル内容を確認
	installedContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("インストールされたファイルの読み取りに失敗: %v", err)
	}
	
	if string(installedContent) != string(templateContent) {
		t.Errorf("インストールされたファイルの内容が一致しません。\n期待値: %s\n実際: %s", 
			string(templateContent), string(installedContent))
	}
}

// Windsurfインストーラーのグローバルインストールテスト
func TestWindsurfInstaller_InstallGlobal(t *testing.T) {
	// テスト環境のセットアップ
	_, templateDir, destDir := setupInstallerTest(t)
	
	// テンプレートファイルを作成
	templateFilePath := filepath.Join(templateDir, "cascade.global.json")
	templateContent := []byte(`{"test": "windsurf-global-config"}`)
	err := os.WriteFile(templateFilePath, templateContent, 0644)
	if err != nil {
		t.Fatalf("テンプレートファイルの作成に失敗: %v", err)
	}
	
	// インストーラーを作成
	installer := NewWindsurfInstaller()
	
	// テスト用にインストール先ディレクトリとテンプレートディレクトリを上書き
	installer.templateDir = templateDir
	installer.globalDestDir = destDir
	
	// グローバルファイルのインストールをテスト
	err = installer.InstallGlobal()
	if err != nil {
		t.Fatalf("グローバルファイルのインストールに失敗: %v", err)
	}
	
	// インストールされたファイルが存在するか確認
	destFile := filepath.Join(destDir, installer.globalFileName)
	exists := fileExists(destFile)
	if !exists {
		t.Errorf("ファイルがインストールされていません: %s", destFile)
	}
	
	// ファイル内容を確認
	installedContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("インストールされたファイルの読み取りに失敗: %v", err)
	}
	
	if string(installedContent) != string(templateContent) {
		t.Errorf("インストールされたファイルの内容が一致しません。\n期待値: %s\n実際: %s", 
			string(templateContent), string(installedContent))
	}
}

// すべてのファイルをインストールするテスト
func TestWindsurfInstaller_InstallAll(t *testing.T) {
	// テスト環境のセットアップ
	_, templateDir, destDir := setupInstallerTest(t)
	
	// テンプレートファイルを作成
	localTemplate := filepath.Join(templateDir, "cascade.local.json")
	localContent := []byte(`{"test": "local-config"}`)
	err := os.WriteFile(localTemplate, localContent, 0644)
	if err != nil {
		t.Fatalf("ローカルテンプレートファイルの作成に失敗: %v", err)
	}
	
	globalTemplate := filepath.Join(templateDir, "cascade.global.json")
	globalContent := []byte(`{"test": "global-config"}`)
	err = os.WriteFile(globalTemplate, globalContent, 0644)
	if err != nil {
		t.Fatalf("グローバルテンプレートファイルの作成に失敗: %v", err)
	}
	
	// インストーラーを作成
	installer := NewWindsurfInstaller()
	installer.templateDir = templateDir
	installer.localDestDir = destDir
	installer.globalDestDir = destDir
	
	// すべてのファイルをインストール
	err = installer.InstallAll()
	if err != nil {
		t.Fatalf("すべてのファイルのインストールに失敗: %v", err)
	}
	
	// ローカルファイルの確認
	localFile := filepath.Join(destDir, installer.localFileName)
	if !fileExists(localFile) {
		t.Errorf("ローカルファイルがインストールされていません: %s", localFile)
	}
	
	// グローバルファイルの確認
	globalFile := filepath.Join(destDir, installer.globalFileName)
	if !fileExists(globalFile) {
		t.Errorf("グローバルファイルがインストールされていません: %s", globalFile)
	}
}
