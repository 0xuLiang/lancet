package fs

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestFileByName(t *testing.T) {
	// 创建一个临时目录
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// 创建一些测试文件
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, file := range files {
		if _, err := os.Create(filepath.Join(dir, file)); err != nil {
			t.Fatal(err)
		}
	}

	// 测试 GetLatestFileByName 函数
	latestFile, err := GetLatestFileByName(path.Join(dir, "*.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if latestFile != filepath.Join(dir, "file3.txt") {
		t.Errorf("expected file3.txt, got %s", latestFile)
	}
}

func TestGetLatestFileByModTime(t *testing.T) {
	// 创建一个临时目录
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// 创建一些测试文件
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, file := range files {
		f, err := os.Create(filepath.Join(dir, file))
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	// 测试 GetLatestFileByModTime 函数
	latestFile, err := GetLatestFileByModTime(filepath.Join(dir, "*.txt"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, filepath.Join(dir, "file3.txt"), latestFile, "expected %s, got %s", filepath.Join(dir, "file3.txt"), latestFile)
}

// CSVRecord 表示 CSV 文件中的一行数据
type CSVRecord struct {
	Key   string
	Value string
}

func TestReadJsonFile(t *testing.T) {
	// 创建一个临时的 JSON 文件
	tempFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 写入一些测试数据到临时文件
	testData := map[string]string{"key": "value"}
	data, _ := json.Marshal(testData)
	if _, err := tempFile.Write(data); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	// 创建一个变量来接收解析的 JSON 数据
	var result map[string]string

	// 使用 ReadJsonFile 函数来读取和解析 JSON 文件的内容
	err = ReadJsonFile(&result, tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 验证 ReadJsonFile 函数的结果
	assert.Equal(t, testData, result)

	// 测试错误情况
	err = ReadJsonFile(&result, "nonexistent.json")
	assert.Error(t, err)
}

func TestReadCSVFile(t *testing.T) {
	// 创建一个临时的 CSV 文件
	tempFile, err := os.CreateTemp("", "*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 创建 CSV writer
	writer := csv.NewWriter(tempFile)

	// 写入表头
	if err := writer.Write([]string{"Key", "Value"}); err != nil {
		t.Fatal(err)
	}

	// 写入一些测试数据到临时文件
	testData := CSVRecord{Key: "key", Value: "value"}
	if err := writer.Write([]string{testData.Key, testData.Value}); err != nil {
		t.Fatal(err)
	}

	// 清空缓存，将所有数据写入文件
	writer.Flush()

	// 关闭文件
	tempFile.Close()

	// 创建一个变量来接收解析的 CSV 数据
	var result []CSVRecord

	// 使用 ReadCSVFile 函数来读取和解析 CSV 文件的内容
	err = ReadCSVFile(&result, tempFile.Name())
	if err != nil {
		t.Fatalf("ReadCSVFile returned error: %v", err)
	}

	// 验证 ReadCSVFile 函数的结果
	assert.Equal(t, []CSVRecord{testData}, result)

	// 测试错误情况
	err = ReadCSVFile(&result, "nonexistent.csv")
	assert.Error(t, err)
}

func TestWriteJsonFile(t *testing.T) {
	// 创建一个临时的 JSON 文件
	tempFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 创建一些测试数据
	testData := map[string]string{"key": "value"}

	// 使用 WriteJsonFile 函数将测试数据写入文件
	err = WriteJsonFile(testData, tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 读取文件的内容
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 解析文件的内容
	var result map[string]string
	err = json.Unmarshal(content, &result)
	if err != nil {
		t.Fatal(err)
	}

	// 验证结果
	assert.Equal(t, testData, result)
}

func TestWriteCSVFile(t *testing.T) {
	// 创建一个临时的 CSV 文件
	tempFile, err := os.CreateTemp("", "*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 创建一些测试数据
	testData := []CSVRecord{{Key: "key", Value: "value"}}

	// 使用 WriteCSVFile 函数将测试数据写入文件
	err = WriteCSVFile(testData, tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 读取文件的内容
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 解析文件的内容
	r := csv.NewReader(strings.NewReader(string(content)))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	// 验证结果
	assert.Equal(t, []string{"Key", "Value"}, records[0])                     // 验证表头
	assert.Equal(t, []string{testData[0].Key, testData[0].Value}, records[1]) // 验证数据
}

func TestReadAndWriteYAMLFile(t *testing.T) {
	// 创建一个临时的 YAML 文件
	tempFile, err := os.CreateTemp("", "*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 创建一些测试数据
	testData := CSVRecord{Key: "key", Value: "value"}

	// 使用 WriteFile 函数将测试数据写入文件
	err = WriteFile(testData, tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 创建一个变量来接收解析的 YAML 数据
	var result CSVRecord

	// 使用 ReadFile 函数来读取和解析 YAML 文件的内容
	err = ReadFile(&result, tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 验证 ReadFile 函数的结果
	assert.Equal(t, testData, result)

	// 测试错误情况
	err = ReadFile(&result, "nonexistent.yaml")
	assert.Error(t, err)
}
