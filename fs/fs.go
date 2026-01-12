package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/0xuLiang/lancet/csv"
	"github.com/gookit/goutil/fsutil"
	"gopkg.in/yaml.v3"
)

// ReadJsonFile 从最新的 JSON 文件中读取数据
func ReadJsonFile(path string, out any) error {
	return ReadFile(path, out, json.Unmarshal)
}

// ReadCSVFile 从最新的 CSV 文件中读取数据
func ReadCSVFile(path string, out any) error {
	return ReadFile(path, out, csv.Unmarshal)
}

// ReadYAMLFile 从最新的 YAML 文件中读取数据
func ReadYAMLFile(path string, out any) error {
	return ReadFile(path, out, yaml.Unmarshal)
}

// WriteJsonFile 将 data 写入到 JSON 文件中，如果 path 中包含 *，则会替换为当前时间戳（格式为 20060102_150405）
func WriteJsonFile(path string, data any) error {
	return WriteFile(path, data, json.Marshal)
}

// WriteCSVFile 将 data 写入到 CSV 文件中，如果 path 中包含 *，则会替换为当前时间戳（格式为 20060102_150405）
func WriteCSVFile(path string, data any) error {
	return WriteFile(path, data, csv.Marshal)
}

// WriteYAMLFile 将 data 写入到 YAML 文件中，如果 path 中包含 *，则会替换为当前时间戳（格式为 20060102_150405）
func WriteYAMLFile(path string, data any) error {
	return WriteFile(path, data, yaml.Marshal)
}

type unmarshal func([]byte, any) error
type marshal func(any) ([]byte, error)

// ReadFile 从最新的文件中读取数据，没有指定 unmarshal 时，会根据后缀名自动选择对应类型的 unmarshal
func ReadFile(path string, out any, unmarshal ...unmarshal) error {
	filename, err := GetLatestFileByName(path)
	if err != nil {
		return fmt.Errorf("get latest file: %w", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if len(unmarshal) == 0 {
		switch ext := filepath.Ext(filename); ext {
		case ".csv":
			unmarshal = append(unmarshal, csv.Unmarshal)
		case ".json":
			unmarshal = append(unmarshal, json.Unmarshal)
		case ".yaml", ".yml":
			unmarshal = append(unmarshal, yaml.Unmarshal)
		default:
			return fmt.Errorf("unsupported file format: %s", ext)
		}
	}

	if err = unmarshal[0](data, out); err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	return nil
}

// WriteFile 将 data 写入到文件中，如果 path 中包含 *，则会替换为当前时间戳（格式为 20060102_150405）
func WriteFile(path string, data any, marshal ...marshal) error {
	if len(marshal) == 0 {
		switch ext := filepath.Ext(path); ext {
		case ".csv":
			marshal = append(marshal, csv.Marshal)
		case ".json":
			marshal = append(marshal, json.Marshal)
		case ".yaml", ".yml":
			marshal = append(marshal, yaml.Marshal)
		default:
			return fmt.Errorf("unsupported file format: %s", ext)
		}
	}

	bs, err := marshal[0](data)
	if err != nil {
		return err
	}

	return SaveFile(path, bs)
}

func SaveFile(path string, data any, optFns ...fsutil.OpenOptionFunc) error {
	return fsutil.SaveFile(TimestampFileName(path), data, optFns...)
}

// TimestampFileName 将 path 中的 * 替换为当前时间戳（格式为 20060102_150405）
func TimestampFileName(path string) string {
	timestamp := time.Now().Format("20060102_150405")
	return strings.Replace(path, "*", timestamp, -1)
}

// GetLatestFileByName 获取最新的文件，基于文件名中的时戳
func GetLatestFileByName(path string) (string, error) {
	matches, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", errors.New("no matching files found")
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i] > matches[j]
	})

	return matches[0], nil
}

// GetLatestFileByModTime 获取最新的文件，基于文件的修改时间
func GetLatestFileByModTime(path string) (string, error) {
	matches, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", errors.New("no matching files found")
	}

	var latestFile string
	var latestTime time.Time
	for _, file := range matches {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return "", err
		}
		if latestFile == "" || fileInfo.ModTime().After(latestTime) {
			latestFile = file
			latestTime = fileInfo.ModTime()
		}
	}

	return latestFile, nil
}
