package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// エンドツーエンドテスト用のヘルパー関数
func buildBinary(t *testing.T, tmpDir string) string {
	// プロジェクトのルートディレクトリを見つける
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("プロジェクトルートの検出に失敗: %v", err)
	}
	
	// バイナリのパス
	binaryPath := filepath.Join(tmpDir, "airules")
	
	// ビルドコマンド
	buildCmd := exec.Command("go", "build", "-o", binaryPath, 
		filepath.Join(projectRoot, "main.go"))
	
	// 作業ディレクトリを設定
	buildCmd.Dir = projectRoot
	
	// ビルド実行
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ビルドに失敗: %v\n%s", err, buildOutput)
	}
	
	return binaryPath
}

// プロジェクトのルートディレクトリを見つける
func findProjectRoot() (string, error) {
	// 現在の作業ディレクトリを取得
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	// go.modファイルを探す
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		
		// 親ディレクトリに移動
		parent := filepath.Dir(dir)
		if parent == dir {
			// ルートに到達
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

// Windsurfコマンドのエンドツーエンドテスト
func TestWindsurfE2E(t *testing.T) {
	// CI環境では実行しない
	if os.Getenv("CI") != "" {
		t.Skip("CI環境ではエンドツーエンドテストをスキップします")
	}
	
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	
	// バイナリをビルド
	binaryPath := buildBinary(t, tmpDir)
	
	// テンプレートディレクトリをセットアップ
	templateDir := filepath.Join(tmpDir, "templates")
	err := os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("テンプレートディレクトリの作成に失敗: %v", err)
	}
	
	// テンプレートファイルを作成
	localTemplateFile := filepath.Join(templateDir, "windsurf.json")
	localTemplateContent := []byte(`{"test": "windsurf-e2e-test"}`)
	err = os.WriteFile(localTemplateFile, localTemplateContent, 0644)
	if err != nil {
		t.Fatalf("テンプレートファイルの作成に失敗: %v", err)
	}
	
	// 環境変数を保存
	oldHome := os.Getenv("HOME")
	oldTemplateDir := os.Getenv("AIRULES_TEMPLATE_DIR")
	
	// テスト用の環境変数を設定
	os.Setenv("HOME", tmpDir)
	os.Setenv("AIRULES_TEMPLATE_DIR", templateDir)
	
	// テスト終了時に環境変数を元に戻す
	defer func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("AIRULES_TEMPLATE_DIR", oldTemplateDir)
	}()
	
	// コマンドを実行
	cmd := exec.Command(binaryPath, "windsurf", "--local")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("コマンド実行に失敗: %v\n%s", err, output)
	}
	
	// インストールされたファイルを確認
	// 注: 実際のパスは実装によって異なる場合があります
	configDir := filepath.Join(tmpDir, ".config", "windsurf")
	configFile := filepath.Join(configDir, "config.json")
	
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("ファイルがインストールされていません: %s", configFile)
		// ディレクトリ構造を確認
		t.Logf("テスト環境のディレクトリ構造:")
		listDirRecursive(t, tmpDir, 0)
	} else if err != nil {
		t.Errorf("ファイル確認中にエラーが発生: %v", err)
	} else {
		t.Logf("ファイルが正常にインストールされました: %s", configFile)
	}
}

// Cursorコマンドのエンドツーエンドテスト
func TestCursorE2E(t *testing.T) {
	// CI環境では実行しない
	if os.Getenv("CI") != "" {
		t.Skip("CI環境ではエンドツーエンドテストをスキップします")
	}
	
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	
	// バイナリをビルド
	binaryPath := buildBinary(t, tmpDir)
	
	// テンプレートディレクトリをセットアップ
	templateDir := filepath.Join(tmpDir, "templates")
	err := os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("テンプレートディレクトリの作成に失敗: %v", err)
	}
	
	// テンプレートファイルを作成
	localTemplateFile := filepath.Join(templateDir, "cursor.json")
	localTemplateContent := []byte(`{"test": "cursor-e2e-test"}`)
	err = os.WriteFile(localTemplateFile, localTemplateContent, 0644)
	if err != nil {
		t.Fatalf("テンプレートファイルの作成に失敗: %v", err)
	}
	
	// 環境変数を保存
	oldHome := os.Getenv("HOME")
	oldTemplateDir := os.Getenv("AIRULES_TEMPLATE_DIR")
	
	// テスト用の環境変数を設定
	os.Setenv("HOME", tmpDir)
	os.Setenv("AIRULES_TEMPLATE_DIR", templateDir)
	
	// テスト終了時に環境変数を元に戻す
	defer func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("AIRULES_TEMPLATE_DIR", oldTemplateDir)
	}()
	
	// コマンドを実行
	cmd := exec.Command(binaryPath, "cursor", "--local")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("コマンド実行に失敗: %v\n%s", err, output)
	}
	
	// インストールされたファイルを確認
	configDir := filepath.Join(tmpDir, ".config", "cursor")
	configFile := filepath.Join(configDir, "config.json")
	
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("ファイルがインストールされていません: %s", configFile)
		// ディレクトリ構造を確認
		t.Logf("テスト環境のディレクトリ構造:")
		listDirRecursive(t, tmpDir, 0)
	} else if err != nil {
		t.Errorf("ファイル確認中にエラーが発生: %v", err)
	} else {
		t.Logf("ファイルが正常にインストールされました: %s", configFile)
	}
}

// ディレクトリ構造を再帰的に表示（デバッグ用）
func listDirRecursive(t *testing.T, dir string, depth int) {
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Logf("%s%s (読み取りエラー: %v)", indent(depth), dir, err)
		return
	}
	
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			t.Logf("%s%s/", indent(depth), file.Name())
			listDirRecursive(t, path, depth+1)
		} else {
			t.Logf("%s%s", indent(depth), file.Name())
		}
	}
}

// インデントを生成
func indent(depth int) string {
	result := ""
	for i := 0; i < depth; i++ {
		result += "  "
	}
	return result
}
