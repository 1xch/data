package data

import "testing"

func TestContainer(t *testing.T) {
	c1 := base
	c1.Retag("TEST_ONE")

	c2 := c1.Clone("a.list", "a.map")
	c2.Retag("TEST_TWO")

	c3 := c1.CloneAs("TEST_THREE")
	c3.Reset()

	t1, t2, t3 := c1.Tag(), c2.Tag(), c3.Tag()
	if t1 != "TEST_ONE" || t2 != "TEST_TWO" || t3 != "TEST_THREE" {
		t.Errorf("error in tagging containers, expected 'TEST_ONE', 'TEST_TWO', 'TEST_THREE' got %s, %s, %s", t1, t2, t3)
	}

	k1, k2, k3 := c1.Keys(), c2.Keys(), c3.Keys()
	l1, l2, l3 := len(k1), len(k2), len(k3)
	if l1 != 11 || l2 != 9 || l3 != 2 {
		t.Errorf("keys lists incorrect expected keys lists of 11,9, 2 and received %d, %d, %d", l1, l2, l3)
	}

	i1, i2, i3 := c1.Get("a.int"), c2.Get("a.int"), c3.Get("a.int")
	i1v, i2v := i1.ToInt(), i2.ToInt()
	if i1v != i2v {
		t.Errorf("returned item values are not equal; %d and %d ", i1v, i2v)
	}
	if i3 != nil {
		t.Errorf("expected nil item, received %v", i3)
	}

	a1, a2 := c1.Match("always"), c2.Match("always")
	la1, la2 := len(a1), len(a2)
	if la1 != 2 || la2 != 2 || la1 != la2 {
		t.Errorf("returned matching items not expected length of value: 1(%v) 2(%v)", a1, a2)
	}

	c3.Set(NewItem("mergable", "mergable"))
	c1.Merge(c3)
	merged := c1.Get("mergable")
	if merged == nil {
		t.Error("unexpected nil item after merging")
	}

	td := c1.TemplateData()
	if v, exists := td["ContainerTag"]; v != "TEST_ONE" || !exists {
		t.Errorf("incorrect template data item: %v", v)
	}

	c3.Clear()
	ks := c3.Keys()
	if len(ks) != 0 {
		t.Errorf("cleared keys length should be zero but was not: existing keys %v", ks)
	}
}
