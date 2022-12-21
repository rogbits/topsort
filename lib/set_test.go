package lib

import (
	"testing"
)

type TestStruct struct {
	A string
	B int
}

func TestNewSet(t *testing.T) {
	s1 := NewSet[string]()
	s1.Add("string")

	s2 := NewSet[int]()
	s2.Add(42)

	s3 := NewSet[TestStruct]()
	s3.Add(TestStruct{"A", 42})
}

func TestSet_Add(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(1)
	s.Add(1)
	if len(s.Map) != 1 {
		t.Fatal("expecting 1")
	}

	s.Add(2)
	s.Add(2)
	if len(s.Map) != 2 {
		t.Fatal("expecting 2")
	}

	s.Add(3)
	s.Add(4)
	if len(s.Map) != 4 {
		t.Fatal("expecting 4")
	}
}

func TestSet_Delete(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	s.Delete(1)
	s.Delete(1)
	if len(s.Map) != 2 {
		t.Fatal("expecting 2")
	}

	s.Delete(2)
	if len(s.Map) != 1 {
		t.Fatal("expecting 1")
	}

	s.Delete(3)
	if len(s.Map) != 0 {
		t.Fatal("expecting 0")
	}
}

func TestSet_Has(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	if !s.Has(1) {
		t.Fatal("expecting 1")
	}
	if !s.Has(2) {
		t.Fatal("expecting 2")
	}
	if !s.Has(3) {
		t.Fatal("expecting 3")
	}

	s.Delete(2)
	if s.Has(2) {
		t.Fatal("expecting 2 removed from set")
	}
}

func TestSet_Iterator(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	for item := range s.Iterator() {
		s.Delete(item)
	}
	if s.Size != 0 {
		t.Fatal("issue during set iteration")
	}
}

func TestSet_GetItems(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)

	items := s.Items()
	if len(items) != 4 {
		t.Fatal("expecting 4")
	}
}
