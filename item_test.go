package data

import "testing"

func TestItem(t *testing.T) {
	if vs := si.ToString(); vs != "string" {
		t.Errorf("string item value is not 'string' it is %v", vs)
	}
	si.SetString("opposite of string")
	if vs2 := si.ToString(); vs2 != "opposite of string" {
		t.Errorf("string item value is not 'opposite of string' it is %v", vs2)
	}
	if k := si.Key(); k != "a.string" {
		t.Errorf("key value is not 'a.string' it is %v", k)
	}
	if u := si.Undotted(); u != "AString" {
		t.Errorf("undotted key value is not 'AString' it is %v", u)
	}
	si.Change("b.string")
	if k := si.Key(); k != "b.string" {
		t.Errorf("key value is not changed to 'b.string' it is %v", k)
	}
	if kv := si.ToKVString(); kv != "b.string:opposite of string" {
		t.Errorf("kv string value is not 'b.string:opposite of string' it is %v", kv)
	}
	a, b := si.ToKV()
	b = b.(string)
	if a != "b.string" || b != "opposite of string" {
		t.Errorf("kv string value is neither 'b.string' nor 'opposite of string' it is %v %v", a, b)
	}
	c := si.Clone("c.string")
	if kv := c.ToKVString(); kv != "c.string:opposite of string" {
		t.Errorf("cloned item kv string value is not 'c.string:opposite of string' it is %v", kv)
	}

	if vb := bi.ToBool(); !vb {
		t.Errorf("bool item value is not 'true' it is %v", vb)
	}
	bi.SetBool(false)
	if vb2 := bi.ToBool(); vb2 {
		t.Errorf("bool item value is not 'false' it is %v", vb2)
	}

	if vi := ii.ToInt(); vi != 9 {
		t.Errorf("int item value is not '9', it is %v", vi)
	}
	ii.SetInt(10)
	if vi2 := ii.ToInt(); vi2 != 10 {
		t.Errorf("int item value is not '10', it is %v", vi2)
	}

	if vf := fi.ToFloat(); vf != 9.9 {
		t.Errorf("float item is not 9.9, it is %v", vf)
	}
	fi.SetFloat(9.8)
	if vf2 := fi.ToFloat(); vf2 != 9.8 {
		t.Errorf("float item is not 9.8, it is %v", vf2)
	}

	if vl := li.ToList(); len(vl) != 7 {
		t.Errorf("list item is not 7 items, it is %v", vl)
	}
	li.SetList("one", "two", "three")
	if vl2 := li.ToList(); len(vl2) != 3 {
		t.Errorf("list item is not 3 items, it is %v", vl2)
	}

	vm := mi.ToMap()
	if vm["c"] != "3" {
		t.Errorf("map item is not properly mapped: %v")
	}
	vm["e"] = "five"
	mi.SetMap(vm)
	if vm2 := mi.ToMap(); vm["e"] != "five" {
		t.Errorf("set map item 'e' is not 'five': %v", vm2)
	}
}
