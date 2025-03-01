package installer

import (
	"os"
	"path/filepath"
	"testing"
)

// ディレクトリ作成機能のテスト
func TestEnsureDirExists(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "test_dir")
	
	// テスト対象の関数を実行
	err := ensureDirExists(testDir)
	if err != nil {
		t.Fatalf("ディレクトリ作成に失敗: %v", err)
	}
	
	// ディレクトリが実際に作成されたか確認
	info, err := os.Stat(testDir)
	if err != nil {
		t.Fatalf("作成したディレクトリを確認できません: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("作成されたパスはディレクトリではありません")
	}
	
	// 既存ディレクトリに対して実行した場合のテスト
	err = ensureDirExists(testDir)
	if err != nil {
		t.Errorf("既存ディレクトリに対して失敗すべきではありません: %v", err)
	}
}

// ファイルの存在確認機能のテスト
func TestFileExists(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.txt")
	nonExistingFile := filepath.Join(tmpDir, "non_existing.txt")
	
	// テスト用のファイルを作成
	err := os.WriteFile(existingFile, []byte("テスト"), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}
	
	// 存在するファイルのテスト
	exists := fileExists(existingFile)
	if !exists {
		t.Errorf("存在するファイルが存在しないと判定されました")
	}
	
	// 存在しないファイルのテスト
	exists = fileExists(nonExistingFile)
	if exists {
		t.Errorf("存在しないファイルが存在すると判定されました")
	}
}

// ファイルコピー機能のテスト
func TestCopyFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "destination.txt")
	
	// テスト用のソースファイルを作成
	testContent := []byte("テストコンテンツ")
	err := os.WriteFile(srcPath, testContent, 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}
	
	// テスト対象の関数を実行
	err = copyFile(srcPath, dstPath)
	if err != nil {
		t.Fatalf("ファイルコピーに失敗: %v", err)
	}
	
	// コピーされたファイルの内容を確認
	copiedContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("コピーしたファイルを読み取れません: %v", err)
	}
	
	if string(copiedContent) != string(testContent) {
		t.Errorf("コピーされた内容が一致しません。\n期待値: %s\n実際: %s", testContent, copiedContent)
	}
}

// バックアップ機能のテスト
func TestBackupFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	originalFile := filepath.Join(tmpDir, "original.txt")
	
	// テスト用のファイルを作成
	originalContent := []byte("オリジナルコンテンツ")
	err := os.WriteFile(originalFile, originalContent, 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}
	
	// テスト対象の関数を実行
	err = backupFile(originalFile)
	if err != nil {
		t.Fatalf("ファイルバックアップに失敗: %v", err)
	}
	
	// バックアップファイルが存在するか確認
	backupPath := originalFile + ".bak"
	exists := fileExists(backupPath)
	if !exists {
		t.Errorf("バックアップファイルが作成されていません: %s", backupPath)
	}
	
	// バックアップファイルの内容を確認
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("バックアップファイルの読み取りに失敗: %v", err)
	}
	
	if string(backupContent) != string(originalContent) {
		t.Errorf("バックアップの内容がオリジナルと一致しません")
	}
}
