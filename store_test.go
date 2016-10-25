package data

import (
	"os"
	"testing"
)

func TestStore(t *testing.T) {
	c1 := base
	rs := c1.ToList("store.retrieval.string")
	s := GetStore(rs[0], rs)
	s.Swap(c1)
	s.Out()

	c2, err := s.In()
	if err != nil {
		t.Error(err)
	}

	lk1, lk2 := len(c1.Keys()), len(c2.Keys())
	if lk1 != lk2 {
		t.Errorf("key length of containers not equal")
	}
	t1, t2 := c1.Tag(), c2.Tag()
	if t1 != t2 {
		t.Errorf("container tags not equal, but should be %s != %s", t1, t2)
	}

	os.Remove(jsonLoc)
}
