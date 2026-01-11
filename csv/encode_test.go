package csv

import (
	"reflect"
	"testing"
)

type Ticket struct {
	Name     string `csv:"name"`
	UserID   string `csv:"user_id"`
	Ticket   int    `csv:"ticket"`
	RecordID string `csv:"record_id"`
	Source   string `csv:"source"`
}

type Simple struct {
	Name string `csv:"name"`
}

// 结构切片
func TestMarshal_1(t *testing.T) {
	tickets := []Ticket{
		{Name: "Alice", UserID: "U001", Ticket: 1, RecordID: "R001", Source: "S001"},
		{Name: "Bob", UserID: "U002", Ticket: 2, RecordID: "R002", Source: "S002"},
	}

	data, err := Marshal(tickets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `name,user_id,ticket,record_id,source
Alice,U001,1,R001,S001
Bob,U002,2,R002,S002
`
	if string(data) != expected {
		t.Errorf("unexpected result: got %v, want %v", string(data), expected)
	}
}

// 结构指针切片
func TestMarshal_2(t *testing.T) {
	tickets := []*Ticket{
		{Name: "Alice", UserID: "U001", Ticket: 1, RecordID: "R001", Source: "S001"},
		{Name: "Bob", UserID: "U002", Ticket: 2, RecordID: "R002", Source: "S002"},
	}

	data, err := Marshal(tickets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `name,user_id,ticket,record_id,source
Alice,U001,1,R001,S001
Bob,U002,2,R002,S002
`
	if string(data) != expected {
		t.Errorf("unexpected result: got %v, want %v", string(data), expected)
	}
}

// 单结构
func TestMarshal_3(t *testing.T) {
	ticket := Ticket{Name: "Alice", UserID: "U001", Ticket: 1, RecordID: "R001", Source: "S001"}

	data, err := Marshal(ticket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `name,user_id,ticket,record_id,source
Alice,U001,1,R001,S001
`
	if string(data) != expected {
		t.Errorf("unexpected result: got %v, want %v", string(data), expected)
	}
}

// 单结构指针
func TestMarshal_4(t *testing.T) {
	ticket := Ticket{Name: "Alice", UserID: "U001", Ticket: 1, RecordID: "R001", Source: "S001"}

	data, err := Marshal(&ticket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `name,user_id,ticket,record_id,source
Alice,U001,1,R001,S001
`
	if string(data) != expected {
		t.Errorf("unexpected result: got %v, want %v", string(data), expected)
	}
}

func TestMarshal_NilPointer(t *testing.T) {
	var ticket *Ticket
	if _, err := Marshal(ticket); err == nil {
		t.Fatalf("expected error for nil pointer")
	}
}

func TestMarshal_SliceWithNilElement(t *testing.T) {
	tickets := []*Ticket{nil}
	if _, err := Marshal(tickets); err == nil {
		t.Fatalf("expected error for nil element")
	}
}

func TestMarshal_PointerToSlice(t *testing.T) {
	tickets := []Ticket{{Name: "Alice", UserID: "U001", Ticket: 1, RecordID: "R001", Source: "S001"}}

	data, err := Marshal(&tickets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `name,user_id,ticket,record_id,source
Alice,U001,1,R001,S001
`
	if string(data) != expected {
		t.Errorf("unexpected result: got %v, want %v", string(data), expected)
	}
}

// 结构切片
func TestUnmarshal_1(t *testing.T) {
	data := []byte(`name,user_id,ticket,record_id,source
Alice,U001,,R001,S001
Alice,U001,1,R001,S001
Bob,U002,2,R002,S002
`)

	var tickets []Ticket
	err := Unmarshal(data, &tickets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []Ticket{
		{Name: "Alice", UserID: "U001", Ticket: 0, RecordID: "R001", Source: "S001"},
		{Name: "Alice", UserID: "U001", Ticket: 1, RecordID: "R001", Source: "S001"},
		{Name: "Bob", UserID: "U002", Ticket: 2, RecordID: "R002", Source: "S002"},
	}
	if !reflect.DeepEqual(tickets, expected) {
		t.Errorf("unexpected result: got %v, want %v", tickets, expected)
	}
}

// 结构指针切片
func TestUnmarshal_2(t *testing.T) {
	data := []byte(`name,user_id,ticket,record_id,source
Alice,U001,1,R001,S001
Bob,U002,2,R002,S002
`)

	var tickets []*Ticket
	err := Unmarshal(data, &tickets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []*Ticket{
		{Name: "Alice", UserID: "U001", Ticket: 1, RecordID: "R001", Source: "S001"},
		{Name: "Bob", UserID: "U002", Ticket: 2, RecordID: "R002", Source: "S002"},
	}
	if !reflect.DeepEqual(tickets, expected) {
		t.Errorf("unexpected result: got %v, want %v", tickets, expected)
	}
}

// 单结构
func TestUnmarshal_3(t *testing.T) {
	data := []byte(`name,user_id,ticket,record_id,source
Alice,U001,1,R001,S001
Bob,U002,2,R002,S002
`)

	var ticket Ticket
	err := Unmarshal(data, &ticket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := Ticket{Name: "Alice", UserID: "U001", Ticket: 1, RecordID: "R001", Source: "S001"}

	if !reflect.DeepEqual(ticket, expected) {
		t.Errorf("unexpected result: got %v, want %v", ticket, expected)
	}
}

func TestUnmarshal_Empty(t *testing.T) {
	var result []Simple
	if err := Unmarshal([]byte(""), &result); err == nil {
		t.Fatalf("expected error for empty data")
	}
}

func TestUnmarshal_ExtraColumns(t *testing.T) {
	data := []byte(`name,extra
Alice,zzz
`)

	var result []Simple
	if err := Unmarshal(data, &result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0].Name != "Alice" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestUnmarshal_HeaderOnlySingleStruct(t *testing.T) {
	data := []byte(`name,extra
`)

	var result Simple
	if err := Unmarshal(data, &result); err == nil {
		t.Fatalf("expected error when no data rows")
	}
}
