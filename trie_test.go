package data

import (
	"crypto/rand"
	"errors"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

const (
	success = true
	failure = false
)

type testData struct {
	key    string
	value  string
	retVal bool
}

func TestTrie_ConstructorOptions(t *testing.T) {
	trie := NewTrie(MaxPrefixPerNode(16), MaxChildrenPerSparseNode(10))

	if trie.maxPrefixPerNode != 16 {
		t.Errorf("Unexpected trie.maxPrefixPerNode value, expected=%v, got=%v",
			16, trie.maxPrefixPerNode)
	}

	if trie.maxChildrenPerSparseNode != 10 {
		t.Errorf("Unexpected trie.maxChildrenPerSparseNode value, expected=%v, got=%v",
			10, trie.maxChildrenPerSparseNode)
	}
}

func TestTrie_GetNonexistentPrefix(t *testing.T) {
	trie := NewTrie()

	d := []testData{
		{key: "aba", value: "0"},
	}

	for _, v := range d {
		trie.set(NewStringItem(v.key, v.value))
	}

	if item := trie.get(Prefix("baa")); item != nil {
		t.Errorf("Unexpected return value, expected=<nil>, got=%v", item)
	}
}

func TestTrie_RandomKitchenSink(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	const count, size = 750000, 16
	b := make([]byte, count+size+1)
	if _, err := rand.Read(b); err != nil {
		t.Fatal("error generating random bytes", err)
	}
	m := make(map[string]string)
	for i := 0; i < count; i++ {
		m[string(b[i:i+size])] = strconv.Itoa(i + size + 1)
	}
	trie := NewTrie()
	getAndDelete := func(k, v string) {
		i := trie.get(Prefix(k))
		si, ok := i.(StringItem)
		if !ok {
			t.Fatalf("unexpected item type, expecting=%v, got=%v", reflect.TypeOf(k), reflect.TypeOf(i))
		}
		if si == nil {
			t.Fatalf("item not found, prefix=%v", []byte(k))
		}
		s := si.ToString()
		if s != v {
			t.Fatalf("unexpected item value, expecting=%s, got=%s", s, v)
		}
		if !trie.Delete(Prefix(k)) {
			t.Fatalf("delete failed, prefix=%v", []byte(k))
		}
		if i = trie.get(Prefix(k)); i != nil {
			t.Fatalf("unexpected item, expecting=<nil>, got=%v", i)
		}
		if trie.Delete(Prefix(k)) {
			t.Fatalf("extra delete succeeded, prefix=%v", []byte(k))
		}
	}
	for k, v := range m {
		ni := NewStringItem(k, v)
		if !trie.put(ni, false) {
			t.Fatalf("insert failed, prefix=%v", []byte(k))
		}
		if byte(k[size/2]) < 128 {
			getAndDelete(k, v)
			delete(m, k)
		}
	}
	for k, v := range m {
		getAndDelete(k, v)
	}
}

func TestTrie_DeleteRoot(t *testing.T) {
	trie := NewTrie()

	v := testData{"aba", "0", success}

	t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
	i := NewStringItem(v.key, v.value)
	if ok := trie.put(i, false); ok != v.retVal {
		t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
	}

	t.Logf("DELETE prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
	if ok := trie.Delete(Prefix(v.key)); ok != v.retVal {
		t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
	}
}

func TestTrie_DeleteAbsentPrefix(t *testing.T) {
	trie := NewTrie()

	v := testData{"a", "0", success}

	t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
	ni := NewStringItem(v.key, v.value)
	if ok := trie.put(ni, false); ok != v.retVal {
		t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
	}

	d := "ab"
	t.Logf("DELETE prefix=%v, success=%v", d, failure)
	if ok := trie.Delete(Prefix(d)); ok != failure {
		t.Errorf("Unexpected return value, expected=%v, got=%v", failure, ok)
	}
	t.Logf("GET prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
	i := trie.get(Prefix(v.key))
	if si, ok := i.(StringItem); ok {
		if si.ToString() != v.value {
			t.Errorf("Unexpected item, expected=%v, got=%v", v.value, i)
		}
	}
}

// overhead is allowed tolerance for Go's runtime/GC to increase the allocated memory
// (to avoid failing tests on insignificant growth amounts)
const overhead = 4000

func TestTrie_InsertDense(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"aba", "0", success},
		{"abb", "1", success},
		{"abc", "2", success},
		{"abd", "3", success},
		{"abe", "4", success},
		{"abf", "5", success},
		{"abg", "6", success},
		{"abh", "7", success},
		{"abi", "8", success},
		{"abj", "9", success},
		{"abk", "0", success},
		{"abl", "1", success},
		{"abm", "2", success},
		{"abn", "3", success},
		{"abo", "4", success},
		{"abp", "5", success},
		{"abq", "6", success},
		{"abr", "7", success},
		{"abs", "8", success},
		{"abt", "9", success},
		{"abu", "0", success},
		{"abv", "1", success},
		{"abw", "2", success},
		{"abx", "3", success},
		{"aby", "4", success},
		{"abz", "5", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}
}

func TestTrie_InsertDensePreceeding(t *testing.T) {
	trie := NewTrie()
	start := byte(70)
	// create a dense node
	for i := byte(0); i <= DefaultMaxChildrenPerSparseNode; i++ {
		k := string([]byte{start + i})
		ni := NewStringItem(k, "")
		if !trie.put(ni, true) {
			t.Errorf("insert failed, prefix=%v", start+i)
		}
	}
	// insert some preceeding keys
	for i := byte(1); i < start; i *= i + 1 {
		k := string([]byte{start - i})
		ni := NewStringItem(k, "")
		if !trie.put(ni, true) {
			t.Errorf("insert failed, prefix=%v", start-i)
		}
	}
}

func TestTrie_InsertDenseDuplicatePrefixes(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"aba", "0", success},
		{"abb", "0", success},
		{"abc", "2", success},
		{"abd", "3", success},
		{"abe", "4", success},
		{"abf", "5", success},
		{"abg", "6", success},
		{"abh", "7", success},
		{"abi", "8", success},
		{"abj", "9", success},
		{"abk", "0", success},
		{"abl", "1", success},
		{"abm", "2", success},
		{"abn", "3", success},
		{"abo", "4", success},
		{"abp", "5", success},
		{"abq", "6", success},
		{"abr", "7", success},
		{"abs", "8", success},
		{"abt", "9", success},
		{"abu", "0", success},
		{"abv", "1", success},
		{"abw", "2", success},
		{"abx", "3", success},
		{"aby", "4", success},
		{"abz", "5", success},
		{"aba", "0", failure},
		{"abb", "1", failure},
		{"abc", "2", failure},
		{"abd", "3", failure},
		{"abe", "4", failure},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}
}

func TestTrie_DeleteDense(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"aba", "0", success},
		{"abb", "1", success},
		{"abc", "2", success},
		{"abd", "3", success},
		{"abe", "4", success},
		{"abf", "5", success},
		{"abg", "6", success},
		{"abh", "7", success},
		{"abi", "8", success},
		{"abj", "9", success},
		{"abk", "0", success},
		{"abl", "1", success},
		{"abm", "2", success},
		{"abn", "3", success},
		{"abo", "4", success},
		{"abp", "5", success},
		{"abq", "6", success},
		{"abr", "7", success},
		{"abs", "8", success},
		{"abt", "9", success},
		{"abu", "0", success},
		{"abv", "1", success},
		{"abw", "2", success},
		{"abx", "3", success},
		{"aby", "4", success},
		{"abz", "5", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	for _, v := range data {
		t.Logf("DELETE word=%v, success=%v", v.key, v.retVal)
		if ok := trie.Delete([]byte(v.key)); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}
}

func TestTrie_DeleteLeakageDense(t *testing.T) {
	trie := NewTrie()

	genTestData := func() *testData {
		// Generate a random hash as a key.
		key := V4Quick()
		return &testData{key: key, value: "v", retVal: success}
	}

	testSize := 100
	data := make([]*testData, 0, testSize)
	for i := 0; i < testSize; i++ {
		data = append(data, genTestData())
	}

	oldBytes := heapAllocatedBytes()

	// repeat insertion/deletion for 10K times to catch possible memory issues
	for i := 0; i < 10000; i++ {
		for _, v := range data {
			ni := NewStringItem(v.key, v.value)
			if ok := trie.put(ni, false); ok != v.retVal {
				t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
			}
		}

		for _, v := range data {
			if ok := trie.Delete([]byte(v.key)); ok != v.retVal {
				t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
			}
		}
	}

	if newBytes := heapAllocatedBytes(); newBytes > oldBytes+overhead {
		//t.Logf("Size=%d, Total=%d, Trie state:\n%s\n", trie.size(), trie.total(), trie.dump())
		t.Errorf("Heap space leak, grew %d bytes (%d to %d)\n", newBytes-oldBytes, oldBytes, newBytes)
	}

	if numChildren := trie.children.length(); numChildren != 0 {
		t.Errorf("Trie is not empty: %v children found", numChildren)
	}
}

func heapAllocatedBytes() uint64 {
	runtime.GC()

	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	return ms.Alloc
}

func TestTrie_InsertDifferentPrefixes(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepaneeeeeeeeeeeeee", "Pepan Zdepan", success},
		{"Honzooooooooooooooo", "Honza Novak", success},
		{"Jenikuuuuuuuuuuuuuu", "Jenik Poustevnicek", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}
}

func TestTrie_InsertDuplicatePrefixes(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepan", "Pepan Zdepan", success},
		{"Pepan", "Pepan Zdepan", failure},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}
}

func TestTrie_InsertVariousPrefixes(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepan", "Pepan Zdepan", success},
		{"Pepin", "Pepin Omacka", success},
		{"Honza", "Honza Novak", success},
		{"Jenik", "Jenik Poustevnicek", success},
		{"Pepan", "Pepan Dupan", failure},
		{"Karel", "Karel Pekar", success},
		{"Jenik", "Jenik Poustevnicek", failure},
		{"Pepanek", "Pepanek Zemlicka", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}
}

func TestTrie_InsertAndMatchPrefix(t *testing.T) {
	trie := NewTrie()
	t.Log("INSERT prefix=by week")
	ni1 := NewStringItem("by week", "2")
	trie.put(ni1, false)
	t.Log("INSERT prefix=by")
	ni2 := NewStringItem("by", "1")
	trie.put(ni2, false)

	if !trie.Match(Prefix("by")) {
		t.Error("MATCH prefix=by, expected=true, got=false")
	}
}

func TestTrie_SetGet(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepan", "Pepan Zdepan", success},
		{"Pepin", "Pepin Omacka", success},
		{"Honza", "Honza Novak", success},
		{"Jenik", "Jenik Poustevnicek", success},
		{"Pepan", "Pepan Dupan", failure},
		{"Karel", "Karel Pekar", success},
		{"Jenik", "Jenik Poustevnicek", failure},
		{"Pepanek", "Pepanek Zemlicka", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	for _, v := range data {
		t.Logf("SET %q to 10", v.key)
		ni := NewStringItem(v.key, "10")
		trie.put(ni, true)
	}

	for _, v := range data {
		value := trie.get(Prefix(v.key))
		t.Logf("GET %q => %v", v.key, value)
		if i, ok := value.(StringItem); ok {
			iv := i.ToString()
			if iv != "10" {
				t.Errorf("Unexpected return value, %s != 10", iv)
			}
		} else {
			t.Errorf("Expected string item, got %v", value)
		}
	}

	if value := trie.get(Prefix("random crap")); value != nil {
		t.Errorf("Unexpected return value, %v != <nil>", value)
	}
}

func TestTrie_Match(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepan", "Pepan Zdepan", success},
		{"Pepin", "Pepin Omacka", success},
		{"Honza", "Honza Novak", success},
		{"Jenik", "Jenik Poustevnicek", success},
		{"Pepan", "Pepan Dupan", failure},
		{"Karel", "Karel Pekar", success},
		{"Jenik", "Jenik Poustevnicek", failure},
		{"Pepanek", "Pepanek Zemlicka", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	for _, v := range data {
		matched := trie.Match(Prefix(v.key))
		t.Logf("MATCH %q => %v", v.key, matched)
		if !matched {
			t.Errorf("Inserted key %q was not matched", v.key)
		}
	}

	if trie.Match(Prefix("random crap")) {
		t.Errorf("Key that was not inserted matched: %q", "random crap")
	}
}

func TestTrie_MatchFalsePositive(t *testing.T) {
	trie := NewTrie()

	if ok := trie.put(NewStringItem("A", "1"), false); !ok {
		t.Fatal("INSERT prefix=A, item=1 not ok")
	}

	resultMatchSubtree := trie.MatchSubtree(Prefix("A extra"))
	resultMatch := trie.Match(Prefix("A extra"))

	if resultMatchSubtree != false {
		t.Error("MatchSubtree returned false positive")
	}

	if resultMatch != false {
		t.Error("Match returned false positive")
	}
}

func TestTrie_MatchSubtree(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepan", "Pepan Zdepan", success},
		{"Pepin", "Pepin Omacka", success},
		{"Honza", "Honza Novak", success},
		{"Jenik", "Jenik Poustevnicek", success},
		{"Pepan", "Pepan Dupan", failure},
		{"Karel", "Karel Pekar", success},
		{"Jenik", "Jenik Poustevnicek", failure},
		{"Pepanek", "Pepanek Zemlicka", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	for _, v := range data {
		key := Prefix(v.key[:3])
		matched := trie.MatchSubtree(key)
		t.Logf("MATCH_SUBTREE %q => %v", key, matched)
		if !matched {
			t.Errorf("Subtree %q was not matched", v.key)
		}
	}
}

func TestTrie_Visit(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepa", "0", success},
		{"Pepa Zdepa", "1", success},
		{"Pepa Kuchar", "2", success},
		{"Honza", "3", success},
		{"Jenik", "4", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Fatalf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	if err := trie.Visit(func(prefix Prefix, item Item) error {
		name := item.Key()
		t.Logf("VISITING prefix=%q, item=%v", prefix, item)
		if !strings.HasPrefix(string(prefix), name) {
			t.Errorf("Unexpected prefix encountered, %q not a prefix of %q", prefix, name)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestTrie_VisitSkipSubtree(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepa", "0", success},
		{"Pepa Zdepa", "1", success},
		{"Pepa Kuchar", "2", success},
		{"Honza", "3", success},
		{"Jenik", "4", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Fatalf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	vf := func(prefix Prefix, item Item) error {
		t.Logf("VISITING prefix=%q, item=%v", prefix, item)
		if i, ok := item.(StringItem); ok {
			if i.ToString() == "0" {
				t.Logf("SKIP %q", prefix)
				return SkipSubtree
			}
		}
		if strings.HasPrefix(string(prefix), "Pepa") {
			t.Errorf("Unexpected prefix encountered, %q", prefix)
		}
		return nil
	}

	if err := trie.Visit(vf); err != nil {
		t.Fatal(err)
	}
}

func TestTrie_VisitReturnError(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepa", "0", success},
		{"Pepa Zdepa", "1", success},
		{"Pepa Kuchar", "2", success},
		{"Honza", "3", success},
		{"Jenik", "4", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		ni := NewStringItem(v.key, v.value)
		if ok := trie.put(ni, false); ok != v.retVal {
			t.Fatalf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	someErr := errors.New("Something exploded")
	vf := func(prefix Prefix, item Item) error {
		t.Logf("VISITING prefix=%q, item=%v", prefix, item)
		if i, ok := item.(StringItem); ok {
			iv := i.ToString()
			if iv == "3" {
				return someErr
			}
			if iv != "3" {
				t.Errorf("Unexpected prefix encountered, %q", prefix)
			}
		}
		return nil
	}
	if err := trie.Visit(vf); err != nil && err != someErr {
		t.Fatal(err)
	}
}

/*
func TestTrie_VisitSubtree(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepa", 0, success},
		{"Pepa Zdepa", 1, success},
		{"Pepa Kuchar", 2, success},
		{"Honza", 3, success},
		{"Jenik", 4, success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		if ok := trie.Insert([]byte(v.key), v.value); ok != v.retVal {
			t.Fatalf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	var counter int
	subtreePrefix := []byte("Pep")
	t.Log("VISIT Pep")
	if err := trie.VisitSubtree(subtreePrefix, func(prefix Prefix, item Item) error {
		t.Logf("VISITING prefix=%q, item=%v", prefix, item)
		if !bytes.HasPrefix(prefix, subtreePrefix) {
			t.Errorf("Unexpected prefix encountered, %q does not extend %q",
				prefix, subtreePrefix)
		}
		if len(prefix) > len(data[item.(int)].key) {
			t.Fatalf("Something is rather fishy here, prefix=%q", prefix)
		}
		counter++
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if counter != 3 {
		t.Error("Unexpected number of nodes visited")
	}
}

func TestTrie_VisitPrefixes(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"P", 0, success},
		{"Pe", 1, success},
		{"Pep", 2, success},
		{"Pepa", 3, success},
		{"Pepa Zdepa", 4, success},
		{"Pepa Kuchar", 5, success},
		{"Honza", 6, success},
		{"Jenik", 7, success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		if ok := trie.Insert([]byte(v.key), v.value); ok != v.retVal {
			t.Fatalf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	var counter int
	word := []byte("Pepa")
	if err := trie.VisitPrefixes(word, func(prefix Prefix, item Item) error {
		t.Logf("VISITING prefix=%q, item=%v", prefix, item)
		if !bytes.HasPrefix(word, prefix) {
			t.Errorf("Unexpected prefix encountered, %q is not a prefix of %q",
				prefix, word)
		}
		counter++
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if counter != 4 {
		t.Error("Unexpected number of nodes visited")
	}
}

func TestParticiaTrie_Delete(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepan", "Pepan Zdepan", success},
		{"Honza", "Honza Novak", success},
		{"Jenik", "Jenik Poustevnicek", success},
	}

	for _, v := range data {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		if ok := trie.Insert([]byte(v.key), v.value); ok != v.retVal {
			t.Fatalf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	for _, v := range data {
		t.Logf("DELETE word=%v, success=%v", v.key, v.retVal)
		if ok := trie.Delete([]byte(v.key)); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}
}

func TestParticiaTrie_DeleteLeakageSparse(t *testing.T) {
	trie := NewTrie()

	data := []testData{
		{"Pepan", "Pepan Zdepan", success},
		{"Honza", "Honza Novak", success},
		{"Jenik", "Jenik Poustevnicek", success},
	}

	oldBytes := heapAllocatedBytes()

	for i := 0; i < 10000; i++ {
		for _, v := range data {
			if ok := trie.Insert([]byte(v.key), v.value); ok != v.retVal {
				t.Fatalf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
			}
		}

		for _, v := range data {
			if ok := trie.Delete([]byte(v.key)); ok != v.retVal {
				t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
			}
		}
	}

	if newBytes := heapAllocatedBytes(); newBytes > oldBytes+overhead {
		t.Logf("Size=%d, Total=%d, Trie state:\n%s\n", trie.size(), trie.total(), trie.dump())
		t.Errorf("Heap space leak, grew %d bytes (from %d to %d)\n", newBytes-oldBytes, oldBytes, newBytes)
	}
}

func TestParticiaTrie_DeleteNonExistent(t *testing.T) {
	trie := NewTrie()

	insertData := []testData{
		{"Pepan", "Pepan Zdepan", success},
		{"Honza", "Honza Novak", success},
		{"Jenik", "Jenik Poustevnicek", success},
	}
	deleteData := []testData{
		{"Pepan", "Pepan Zdepan", success},
		{"Honza", "Honza Novak", success},
		{"Pepan", "Pepan Zdepan", failure},
		{"Jenik", "Jenik Poustevnicek", success},
		{"Honza", "Honza Novak", failure},
	}

	for _, v := range insertData {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		if ok := trie.Insert([]byte(v.key), v.value); ok != v.retVal {
			t.Fatalf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	for _, v := range deleteData {
		t.Logf("DELETE word=%v, success=%v", v.key, v.retVal)
		if ok := trie.Delete([]byte(v.key)); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}
}

func TestParticiaTrie_DeleteSubtree(t *testing.T) {
	trie := NewTrie()

	insertData := []testData{
		{"P", 0, success},
		{"Pe", 1, success},
		{"Pep", 2, success},
		{"Pepa", 3, success},
		{"Pepa Zdepa", 4, success},
		{"Pepa Kuchar", 5, success},
		{"Honza", 6, success},
		{"Jenik", 7, success},
	}
	deleteData := []testData{
		{"Pe", -1, success},
		{"Pe", -1, failure},
		{"Honzik", -1, failure},
		{"Honza", -1, success},
		{"Honza", -1, failure},
		{"Pep", -1, failure},
		{"P", -1, success},
		{"Nobody", -1, failure},
		{"", -1, success},
	}

	for _, v := range insertData {
		t.Logf("INSERT prefix=%v, item=%v, success=%v", v.key, v.value, v.retVal)
		if ok := trie.Insert([]byte(v.key), v.value); ok != v.retVal {
			t.Fatalf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}

	for _, v := range deleteData {
		t.Logf("DELETE_SUBTREE prefix=%v, success=%v", v.key, v.retVal)
		if ok := trie.DeleteSubtree([]byte(v.key)); ok != v.retVal {
			t.Errorf("Unexpected return value, expected=%v, got=%v", v.retVal, ok)
		}
	}
}

func TestTrie_compact(t *testing.T) {
	trie := NewTrie()

	trie.Insert(Prefix("a"), 0)
	trie.Insert(Prefix("ab"), 0)
	trie.Insert(Prefix("abc"), 0)
	trie.Insert(Prefix("abcd"), 0)
	trie.Insert(Prefix("abcde"), 0)
	trie.Insert(Prefix("abcdef"), 0)
	trie.Insert(Prefix("abcdefg"), 0)
	trie.Insert(Prefix("abcdefgi"), 0)
	trie.Insert(Prefix("abcdefgij"), 0)
	trie.Insert(Prefix("abcdefgijk"), 0)

	trie.Delete(Prefix("abcdef"))
	trie.Delete(Prefix("abcde"))
	trie.Delete(Prefix("abcdefg"))

	trie.Delete(Prefix("a"))
	trie.Delete(Prefix("abc"))
	trie.Delete(Prefix("ab"))

	trie.Visit(func(prefix Prefix, item Item) error {
		// 97 ~~ 'a',
		for ch := byte(97); ch <= 107; ch++ {
			if c := bytes.Count(prefix, []byte{ch}); c > 1 {
				t.Errorf("%q appeared in %q %v times", ch, prefix, c)
			}
		}
		return nil
	})
}

func TestTrie_longestCommonPrefixLenght(t *testing.T) {
	trie := NewTrie()
	trie.prefix = []byte("1234567890")

	switch {
	case trie.longestCommonPrefixLength([]byte("")) != 0:
		t.Fail()
	case trie.longestCommonPrefixLength([]byte("12345")) != 5:
		t.Fail()
	case trie.longestCommonPrefixLength([]byte("123789")) != 3:
		t.Fail()
	case trie.longestCommonPrefixLength([]byte("12345678901")) != 10:
		t.Fail()
	}
}
*/
