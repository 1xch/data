package data

import (
	"os"
	"testing"
)

func TestStore(t *testing.T) {
	c1 := base
	rs := c1.ToStrings("store.retrieval.string")
	s1 := GetStore(rs[0], rs)
	s1.Swap(c1)
	s1.Out()

	s2 := GetStore("yaml", rs)
	s2.Swap(c1)
	s2.Out()

	c2, err := s1.In()
	if err != nil {
		t.Error(err)
	}

	c3, err := s2.In()
	if err != nil {
		t.Error(err)
	}

	lk1, lk2, lk3 := len(c1.Keys()), len(c2.Keys()), len(c3.Keys())
	if lk1 != lk2 || lk2 != lk3 || lk3 != lk1 {
		t.Errorf("key length of containers not equal")
	}
	t1, t2, t3 := c1.Tag(), c2.Tag(), c3.Tag()
	if t1 != t2 || t2 != t3 || t3 != t1 {
		t.Errorf("container tags not equal, but should be %s != %s != %s", t1, t2, t3)
	}

	os.Remove(jsonLoc)
	os.Remove(yamlLoc)
}

type storeTestFunc func(*testing.T, Store, *Vector)

func storeTest(t *testing.T, trs []string, removal string, fn ...storeTestFunc) {
	c := base.Clone()
	c.Set(
		NewStringsItem("store.retrieval.string", trs...),
	)
	crs := c.ToStrings("store.retrieval.string")
	s := GetStore(crs[0], crs)

	for _, f := range fn {
		f(t, s, c)
	}

	if removal != "" {
		os.Remove(removal)
	}
}

func TestOutStore(t *testing.T) {
	loc := "./out.txt"
	f, err := os.Create(loc)
	if err != nil {
		t.Error(err)
	}
	o := &StoreMaker{
		"OUT", OutStore(f),
	}
	SetStore(o)

	trs := []string{"OUT", currentDir, "out.txt"}

	tf := func(t *testing.T, s Store, c *Vector) {
		s.Swap(c)
		if _, err := s.Out(); err != nil {
			t.Error(err)
		}
		if _, err := s.In(); err == nil {
			t.Error("Error is nil but should not be FunctionNotImplementedError")
		}
	}

	storeTest(t, trs, loc, tf)
}

func containerFromStoreTest(t *testing.T, s Store, c *Vector) {
	c1 := c
	s.Swap(c1)

	//Writer
	//wb := make([]byte, 500)
	//_, _ = s.Write(wb)
	//spew.Dump(wb)

	if _, err := s.Out(); err != nil {
		t.Error(err)
	}

	//Reader
	//rb := make([]byte, 500)
	//_, _ = s.Read(rb)
	//spew.Dump(rb)

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
}

func TestJsonStore(t *testing.T) {
	storeTest(t, rs, jsonLoc, containerFromStoreTest)
}

func TestJsonfStore(t *testing.T) {
	trs := *&rs
	trs[0] = "jsonf"
	storeTest(t, rs, jsonLoc, containerFromStoreTest)
}

func TestYamlStore(t *testing.T) {
	trs := *&rs
	trs[0] = "yaml"
	storeTest(t, rs, yamlLoc, containerFromStoreTest)
}

type testStore struct {
	c *Vector
}

func (s *testStore) Read([]byte) (int, error) {
	return 7, nil
}

func (s *testStore) In() (*Vector, error) {
	return s.c, nil
}

func (s *testStore) Swap(c *Vector) {
	s.c = c
}

func (s *testStore) Out() (*Vector, error) {
	return s.c, nil
}

func (s *testStore) Write([]byte) (int, error) {
	return 77, nil
}

func customStore([]string) Store {
	return &testStore{
		New("CUSTOM"),
	}
}

func TestCustomStore(t *testing.T) {
	mk := &StoreMaker{
		"CUSTOM", customStore,
	}
	SetStore(mk)
	s := GetStore("CUSTOM", []string{"custom", "string"})

	if n, err := s.Read([]byte{}); n != 7 || err != nil {
		t.Errorf("custom store read value is not 7: it is %d with error %s", n, err.Error())
	}

	c1, err := s.In()
	if err != nil {
		t.Error(err)
	}
	if tag := c1.ToString("vector.tag"); tag != "CUSTOM" {
		t.Errorf("Custom store tag is not 'CUSTOM', it is %s", tag)
	}

	c2 := New("CUSTOM#2")
	s.Swap(c2)
	c3, err := s.Out()
	if err != nil {
		t.Error(err)
	}
	t1, t2 := c2.Tag(), c3.Tag()
	if t1 != t2 {
		t.Errorf("returned custom store vector tags not equal, but should be %s != %s", t1, t2)
	}

	if n, err := s.Write([]byte{}); n != 77 || err != nil {
		t.Errorf("custom store read value is not 77: it is %d with error %s", n, err.Error())
	}
}
