package data

import (
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestItem(t *testing.T) {
	i1 := si.Clone()
	if k := i1.Key(); k != "a.string" {
		t.Errorf("key value is not 'a.string' it is %v", k)
	}
	i1.NewKey("b.string")
	if k := i1.Key(); k != "b.string" {
		t.Errorf("key value is not changed to 'b.string' it is %v", k)
	}
	if u := i1.KeyUndotted(); u != "BString" {
		t.Errorf("undotted key value is not 'BString' it is %v", u)
	}

	i2 := si.Clone()

	p1, p2 := i1.Provided(), i2.Provided()
	if p1 != p2 {
		t.Errorf("Provided values are not equal: %v - %v", p1, p2)
	}

	b1, err := i1.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	b2, err := yaml.Marshal(i2)
	if err != nil {
		t.Error(err)
	}

	i3, i4 := NewStringItem("", ""), NewStringItem("", "")

	err = i3.UnmarshalJSON(b1)
	if err != nil {
		t.Error(err)
	}

	err = yaml.Unmarshal(b2, i4)
	if err != nil {
		t.Error(err)
	}

	if i3.Key() != "b.string" || i4.Key() != "a.string" {
		t.Errorf("unmarshaled item keys are incorrect: %v, %v", i3, i4)
	}

	if i3.ToString() != i4.ToString() {
		t.Errorf("unmarshaled item values are incorrect: %v, %v", i3, i4)
	}
}

func TestStringItem(t *testing.T) {
	i, ok := si.(StringItem)
	if !ok {
		t.Errorf("item is not StringItem %v", si)
	}
	if ok {
		if vs := i.ToString(); vs != "string" {
			t.Errorf("string item value is not 'string' it is %v", vs)
		}
		i.SetString("opposite of string")
		if vs2 := i.ToString(); vs2 != "opposite of string" {
			t.Errorf("string item value is not 'opposite of string' it is %v", vs2)
		}
	}
}

func TestStringsItem(t *testing.T) {
	i, ok := ssi.(StringsItem)
	if !ok {
		t.Errorf("item is not StringsItem %v", ssi)
	}
	if ok {
		vl1 := i.ToStrings()
		if len(vl1) != 3 || vl1[2] != "c" {
			t.Errorf("strings item is not ['a', 'b', 'c'], it is %v", vl1)
		}
		i.SetStrings("one", "two", "three", "four")
		vl2 := i.ToStrings()
		if len(vl2) != 4 || vl2[3] != "four" {
			t.Errorf("strings item is not ['one', 'two', 'three', 'four'], it is %v", vl2)
		}
	}
}

func TestBoolItem(t *testing.T) {
	i, ok := bi.(BoolItem)
	if !ok {
		t.Errorf("item is not BoolItem %v", bi)
	}
	if ok {
		if v1 := i.ToBool(); v1 {
			t.Errorf("bool item value is not 'false' it is %v", v1)
		}
		i.SetBool(true)
		if v2 := i.ToBool(); !v2 {
			t.Errorf("bool item value is not 'true' it is %v", v2)
		}

	}
}

func TestIntItem(t *testing.T) {
	i, ok := ii.(IntItem)
	if !ok {
		t.Errorf("item is not IntItem %v", ii)
	}
	if ok {
		if v1 := i.ToInt(); v1 != 9 {
			t.Errorf("int item value is not '9', it is %v", v1)
		}
		i.SetInt(10)
		if v2 := i.ToInt(); v2 != 10 {
			t.Errorf("int item value is not '10', it is %v", v2)
		}
	}
}

func TestInt64Item(t *testing.T) {
	i, ok := ii64.(Int64Item)
	if !ok {
		t.Errorf("item is not Int64Item %v", ii64)
	}
	if ok {
		if v1 := i.ToInt64(); v1 != 999 {
			t.Errorf("int64 item value is not '999', it is %v", v1)
		}
		i.SetInt64(100)
		if v2 := i.ToInt64(); v2 != 100 {
			t.Errorf("int64 item value is not '100', it is %v", v2)
		}
	}
}

func TestUintItem(t *testing.T) {
	i, ok := ui.(UintItem)
	if !ok {
		t.Errorf("item is not UintItem %v", ui)
	}
	if ok {
		if v1 := i.ToUint(); v1 != 1 {
			t.Errorf("uint item value is not '1', it is %v", v1)
		}
		i.SetUint(100)
		if v2 := i.ToUint(); v2 != 100 {
			t.Errorf("uint item value is not '100', it is %v", v2)
		}
	}
}

func TestUint64Item(t *testing.T) {
	i, ok := ui64.(Uint64Item)
	if !ok {
		t.Errorf("item is not Uint64Item %v", ui64)
	}
	if ok {
		if v1 := i.ToUint64(); v1 != 111 {
			t.Errorf("uint64 item value is not '111', it is %v", v1)
		}
		i.SetUint64(100)
		if v2 := i.ToUint64(); v2 != 100 {
			t.Errorf("int item value is not '100', it is %v", v2)
		}
	}
}

func TestFloat64Item(t *testing.T) {
	i, ok := fi.(Float64Item)
	if !ok {
		t.Errorf("item is not FloatItem %v", fi)
	}
	if ok {
		if v1 := i.ToFloat64(); v1 != 9.9 {
			t.Errorf("float item is not 9.9, it is %v", v1)
		}
		i.SetFloat(9.8)
		if v2 := i.ToFloat64(); v2 != 9.8 {
			t.Errorf("float item is not 9.8, it is %v", v2)
		}
	}
}

func TestVectorItem(t *testing.T) {
	i, ok := vi.(VectorItem)
	if !ok {
		t.Errorf("item is not VectorItem %v", vi)
	}
	if ok {
		v1 := i.ToVector()
		ty := reflect.TypeOf(v1)
		if ty.String() != "*data.Vector" {
			t.Errorf("item value is not *data.Vector: %v", ty)
		}
		ci := v1.ToString("vector.1")
		if ci != "ONE" {
			t.Errorf("multi item container item is not 'ONE': %s", ci)
		}
	}
}
