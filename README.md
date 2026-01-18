# lancet

Go 工具集：
- `cond`: 类似三元表达式的条件辅助函数。
- `csv`: 在结构体与 CSV 数据之间进行编解码。
- `fs`: 处理 JSON/CSV/YAML 文件的读写，支持按时间戳输出和自动选择最新文件。

## 安装

```bash
go get github.com/0xuLiang/lancet
```

> 需要 Go 1.25 及以上。

## 快速开始

### 条件工具 `cond`

```go
package main

import (
    "fmt"

    "github.com/0xuLiang/lancet/cond"
)

func main() {
    // Val 像三元表达式一样返回值
    fmt.Println(cond.Val(true, "A", "B")) // 输出: A

    // Fun 延迟执行的分支
    n := cond.Fun(true, func() int { return 1 }, func() int { return 2 })
    fmt.Println(n) // 输出: 1

    // Or 返回第一个非零值（零值则继续向后）
    fmt.Println(cond.Or("", "", "go", "lang")) // 输出: go
}
```

### CSV 编解码 `csv`

```go
package main

import (
    "fmt"

    lancetcsv "github.com/0xuLiang/lancet/csv"
)

type Meta struct {
    Source string `csv:"source"`
}

// 支持匿名（嵌入）字段，会被展平
type Ticket struct {
    Meta
    Name     string `csv:"name"`
    UserID   string `csv:"user_id"`
    Ticket   int    `csv:"ticket"`
    RecordID string `csv:"record_id"`
    Extra    string `csv:"extra,omitempty"` // omitempty: 所有记录均为零值时会整列省略
}

func main() {
    // Marshal: 结构体/切片 -> CSV 字节
    rows := []Ticket{{Meta: Meta{Source: "S001"}, Name: "Alice", UserID: "U001", Ticket: 1, RecordID: "R001"}}
    data, _ := lancetcsv.Marshal(rows)
    fmt.Printf("%s\n", data)

    // Unmarshal: CSV 字节 -> 结构体切片
    var decoded []Ticket
    _ = lancetcsv.Unmarshal(data, &decoded)
    fmt.Println(decoded[0].Name) // Alice
}
```

支持：结构体 / 结构体指针 / 结构体切片 / 结构体指针切片；
- 匿名（嵌入）结构体会被展平，包含指针形式的嵌入；
- `csv:"name,omitempty"`：若所有记录该字段均为零值则整列省略，出现任意非零值即输出该列；
- 时间类型使用 `time.Time` 的文本编解码。

### 文件工具 `fs`

- `ReadJsonFile`/`ReadCSVFile`/`ReadYAMLFile`：读取与路径匹配的最新文件（按文件名排序）并反序列化。
- `WriteJsonFile`/`WriteCSVFile`/`WriteYAMLFile`：写入文件；路径中包含 `*` 时自动替换为时间戳（格式 `20060102_150405`）。
- `GetLatestFileByName` 与 `GetLatestFileByModTime`：按文件名或修改时间获取最新文件。

```go
package main

import (
    "fmt"

    "github.com/0xuLiang/lancet/fs"
)

type Item struct {
    Key   string `csv:"key"`
    Value string `csv:"value"`
}

func main() {
    // 按时间戳写入 CSV
    _ = fs.WriteCSVFile("./out/items-*.csv", []Item{{Key: "k1", Value: "v1"}})

    // 读取匹配的最新 CSV
    var items []Item
    _ = fs.ReadCSVFile("./out/items-*.csv", &items)
    fmt.Println(len(items))
}
```

## 测试

运行全部测试：

```bash
go test ./...
```

## 许可

本项目采用 MIT License。参见 `LICENSE`。
