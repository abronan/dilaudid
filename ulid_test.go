package dilaudid

import (
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

func TestNewRandom(t *testing.T) {
	ulid := NewRandom()
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
	actual2 := ulid.random
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

func BenchmarkULID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewRandom()
	}
}

func BenchmarkEncodedULID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewRandom().String()
	}
}

func BenchmarkSingleEncodedULID(b *testing.B) {
	u := NewRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.String()
	}
}
