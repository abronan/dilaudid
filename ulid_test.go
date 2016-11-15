package dilaudid

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ulid := New(time.Unix(1469918176, 385000000), [10]byte{})
	et := ulid.String()
	if et != "01ARYZ6S410000000000000000" {
		t.Fatalf("expected '01ARYZ6S4100000000000000', got %v", et)
	}
}

func TestUnique(t *testing.T) {
	ulid := Unique()
	et := ulid.String()
	if len(et) != 26 {
		t.Fatalf("expected 26 characters, got %q", et)
	}
	for i := 0; i < 26; i++ {
		if et[i] == 0 {
			t.Fatalf("zero character found at offset %v: %q", i, et)
		}
	}
}

func TestDecode(t *testing.T) {
	ulid, err := Decode("01B1JK8PG4Y5Z1ED14CACZHRBM")
	if err != nil {
		t.Fatalf("error %v", err)
	}

	expected := int64(1479166679556000000)
	actual := ulid.time.UnixNano()
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

	expected2 := [10]byte{241, 126, 23, 52, 36, 98, 153, 248, 225, 116}
	actual2 := ulid.nonce
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected2, actual2)
	}
}

func TestDecodeBad1(t *testing.T) {
	_, err := Decode("01B1JK8PG4Y5Z1ED14CACZHRBMX")
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func TestDecodeBad2(t *testing.T) {
	_, err := Decode("01B1JK8PG4Y5Z1ED14CACZHRBL")
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func TestDecodeBad3(t *testing.T) {
	_, err := Decode("01B1JK8PG4Y5Z1ED14CACZHRB")
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func TestMarshalJSON(t *testing.T) {
	expected := []byte("\"01ARYZ6S410000000000000000\"")
	ulid := New(time.Unix(1469918176, 385000000), [10]byte{})
	marshalled, err := json.Marshal(&ulid)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if !reflect.DeepEqual(marshalled, expected) {
		t.Fatalf("expected %v, got %v", expected, marshalled)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var ulid ULID
	expected := New(time.Unix(1469918176, 385000000), [10]byte{247, 71, 242, 197, 159, 191, 110, 22, 115, 121})
	bytes := []byte("\"01ARYZ6S41YX3Z5HCZQXQ1CWVS\"")
	err := json.Unmarshal(bytes, &ulid)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if ulid.time != expected.time {
		t.Fatalf("expected %v, got %v", expected.time, ulid.time)
	}
	if ulid.nonce != expected.nonce {
		t.Fatalf("expected %v, got %v", expected.nonce, ulid.nonce)
	}
}

func BenchmarkULID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Unique()
	}
}

func BenchmarkEncodedULID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Unique().String()
	}
}

func BenchmarkSingleEncodedULID(b *testing.B) {
	u := Unique()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.String()
	}
}
