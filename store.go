package data

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Store interface {
	Read([]byte) (int, error)
	In() (*Vector, error)
	Swap(*Vector)
	Out() (*Vector, error)
	Write([]byte) (int, error)
}

type ReadFunc func(*Retriever) (io.ReadCloser, int64, error)

type InFunc func(string, int64, io.ReadCloser) (*Vector, error)

type OutFunc func(*Vector, io.WriteCloser) ([]string, error)

type WriteFunc func(*Vector) (io.WriteCloser, error)

type store struct {
	*Retriever
	c   *Vector
	rfn ReadFunc
	ifn InFunc
	ofn OutFunc
	wfn WriteFunc
}

func NewStore(rfn ReadFunc, ifn InFunc, ofn OutFunc, wfn WriteFunc, rs ...string) Store {
	return &store{
		Retriever: NewRetriever(rs...),
		rfn:       rfn,
		ifn:       ifn,
		ofn:       ofn,
		wfn:       wfn,
	}
}

type Retriever struct {
	v []string
}

func NewRetriever(with ...string) *Retriever {
	return &Retriever{with}
}

func (r *Retriever) SetRetrieval(s []string) {
	r.v = s
}

func (r *Retriever) Retrieval() []string {
	return r.v
}

func (r *Retriever) RetrievalString() string {
	return strings.Join(r.v, ":")
}

var MalformedRetrievalStringError = Drror("%s is malformed: %s").Out

func (s *store) Read(p []byte) (int, error) {
	r, _, err := s.rfn(s.Retriever)
	if err != nil {
		return 0, err
	}
	i, err := r.Read(p)
	r.Close()
	return i, err
}

func (s *store) In() (*Vector, error) {
	r, n, err := s.rfn(s.Retriever)
	if err != nil {
		return nil, err
	}
	c, err := s.ifn(s.RetrievalString(), n, r)
	if err != nil {
		return nil, err
	}
	s.Swap(c)
	return c, nil
}

func (s *store) Swap(c *Vector) {
	s.c = c
	s.SetRetrieval(c.ToStrings("store.retrival.string"))
}

func (s *store) Out() (*Vector, error) {
	w, err := s.wfn(s.c)
	if err != nil {
		return nil, err
	}
	r, err := s.ofn(s.c, w)
	if err != nil {
		return nil, err
	}
	s.SetRetrieval(r)
	return s.c, nil
}

func (s *store) Write(p []byte) (int, error) {
	w, err := s.wfn(s.c)
	if err != nil {
		return 0, err
	}
	return w.Write(p)
}

type StoreMaker struct {
	Key string
	Fn  StoreFn
}

type StoreFn func([]string) Store

type Stores map[string]*StoreMaker

var AvailableStores Stores

func SetStore(fs ...*StoreMaker) {
	AvailableStores.Set(fs...)
}

func (s Stores) Set(fs ...*StoreMaker) {
	for _, v := range fs {
		s[v.Key] = v
	}
}

func GetStore(k string, r []string) Store {
	return AvailableStores.Get(k, r)
}

func (s Stores) Get(k string, r []string) Store {
	if f, ok := s[k]; ok {
		return f.Fn(r)
	}
	return nil
}

func init() {
	AvailableStores = make(Stores)
	AvailableStores.Set(StdoutStore, YamlStore, JsonStore, JsonFStore)
}

var FunctionNotImplemented = Drror("%s function not implemented for the %s store.").Out

var StdoutStore = &StoreMaker{"STDOUT", OutStore(os.Stdout)}

func OutStore(out *os.File) StoreFn {
	return func([]string) Store {
		rs := []string{"default", "STDOUT"}
		return NewStore(
			func(rt *Retriever) (io.ReadCloser, int64, error) {
				return nil, 0, FunctionNotImplemented("Read Function", "STDOUT")
			},
			func(r string, n int64, rr io.ReadCloser) (*Vector, error) {
				return nil, FunctionNotImplemented("In Function", "STDOUT")
			},
			func(c *Vector, w io.WriteCloser) ([]string, error) {
				b, err := c.MarshalJSON()
				if err != nil {
					return nil, err
				}
				_, err = w.Write(b)
				w.Close()
				return rs, err
			},
			func(c *Vector) (io.WriteCloser, error) {
				return out, nil
			},
			rs...,
		)
	}
}

var YamlStore = &StoreMaker{"yaml", yamlStore}

func yamlStore(rs []string) Store {
	return NewStore(
		readCloserFrom("yaml"),
		func(r string, n int64, rr io.ReadCloser) (*Vector, error) {
			c := New("")
			b := make([]byte, n)
			_, err := rr.Read(b)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(b[:n], &c)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		func(c *Vector, w io.WriteCloser) ([]string, error) {
			retrieval := c.ToStrings("store.retrieval.string")
			y, err := yaml.Marshal(&c)
			if err != nil {
				return nil, err
			}
			_, err = w.Write(y)
			w.Close()
			return retrieval, err
		},
		writeCloserFrom("yaml"),
		rs...,
	)
}

var (
	JsonStore  = &StoreMaker{"json", JsonStorer(regular)}
	JsonFStore = &StoreMaker{"jsonf", JsonStorer(indented)}
)

type jsonMarshaler func(*Vector) ([]byte, error)

func regular(c *Vector) ([]byte, error) {
	return json.Marshal(&c)
}

func indented(c *Vector) ([]byte, error) {
	return json.MarshalIndent(&c, "", "    ")
}

func JsonStorer(jm jsonMarshaler) StoreFn {
	return func(rs []string) Store {
		return NewStore(
			readCloserFrom("json"),
			func(r string, n int64, rr io.ReadCloser) (*Vector, error) {
				c := New("")
				b := make([]byte, n)
				_, err := rr.Read(b)
				if err != nil {
					return nil, err
				}
				err = c.UnmarshalJSON(b[:n])
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			func(c *Vector, w io.WriteCloser) ([]string, error) {
				retrieval := c.ToStrings("store.retrieval.string")
				j, err := jm(c)
				if err != nil {
					return nil, err
				}
				_, err = w.Write(j)
				w.Close()
				return retrieval, err
			},
			writeCloserFrom("json"),
			rs...,
		)
	}
}

func insufficient(s []string, i int) error {
	if len(s) < i {
		return MalformedRetrievalStringError(s, "expected length equal to or greater than three")
	}
	return nil
}

func writeCloserFrom(ext string) WriteFunc {
	return func(c *Vector) (io.WriteCloser, error) {
		rs := c.ToStrings("store.retrieval.string")
		if err := insufficient(rs, 3); err != nil {
			return nil, err
		}
		loc, fil := rs[1], rs[2]
		p := filepath.Join(loc, fmt.Sprintf("%s.%s", fil, ext))
		fl, err := Open(p)
		if err != nil {
			return nil, err
		}
		fl.Truncate(0)
		return fl, nil
	}
}

var ReaderRetrievalError = Drror("unable to find readcloser: %s").Out

func readCloserFrom(ext string) ReadFunc {
	return func(rt *Retriever) (io.ReadCloser, int64, error) {
		rs := rt.Retrieval()
		if err := insufficient(rs, 3); err != nil {
			return nil, 0, err
		}
		dir, file := rs[1], rs[2]
		fileName := strings.Join([]string{file, ext}, ".")
		info, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, 0, err
		}
		for _, f := range info {
			if !f.IsDir() {
				fn := f.Name()
				if fn == fileName {
					p := filepath.Join(dir, fn)
					fl, err := Open(p)
					if err != nil {
						return nil, 0, err
					}
					var n int64
					if fi, err := fl.Stat(); err == nil {
						if size := fi.Size(); size < 1e9 {
							n = size
						}
					}
					return fl, n, nil
				}
			}
		}
		return nil, 0, ReaderRetrievalError("no suitable path from %v", rs)
	}
}

var openError = Drror("unable to find or open file %s, provided %s").Out

func Exist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModeDir|0755)
	}
}

func Open(path string) (*os.File, error) {
	p := filepath.Clean(path)
	dir, name := filepath.Split(p)
	var fp string
	var err error
	switch dir {
	case "":
		fp, err = filepath.Abs(name)
	default:
		Exist(dir)
		fp, err = filepath.Abs(p)
	}

	if err != nil {
		return nil, err
	}

	if file, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0660); err == nil {
		return file, nil
	}

	return nil, openError(fp, path)
}
