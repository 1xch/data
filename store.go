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
	Out() error
	Write([]byte) (int, error)
	In() (*Container, error)
	Read([]byte) (int, error)
	Swap(*Container)
}

type store struct {
	*retriever
	c   *Container
	ofn outFunc
	wfn writeFunc
	ifn inFunc
	rfn readFunc
}

type outFunc func(*Container, io.WriteCloser) ([]string, error)

type writeFunc func(*Container) (io.WriteCloser, error)

type inFunc func(string, int64, io.ReadCloser) (*Container, error)

type readFunc func(*retriever) (io.ReadCloser, int64, error)

type retriever struct {
	v []string
}

func newRetriever(with ...string) *retriever {
	return &retriever{with}
}

func (r *retriever) SetString(s []string) {
	r.v = s
}

func (r *retriever) Retrieval() []string {
	return r.v
}

func (r *retriever) RetrievalString() string {
	return strings.Join(r.v, ":")
}

var MalformedRetrievalStringError = Drror("%s is malformed: %s").Out

func (s *store) Swap(c *Container) {
	s.c = c
	s.retriever.SetString(c.ToList("store.retrival.string"))
}

func (s *store) Out() error {
	w, err := s.wfn(s.c)
	if err != nil {
		return err
	}
	r, err := s.ofn(s.c, w)
	if err != nil {
		return err
	}
	s.SetString(r)
	return nil
}

func (s *store) Write(p []byte) (int, error) {
	w, err := s.wfn(s.c)
	if err != nil {
		return 0, err
	}
	return w.Write(p)
}

func (s *store) In() (*Container, error) {
	r, n, err := s.rfn(s.retriever)
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

func (s *store) Read(p []byte) (int, error) {
	r, _, err := s.rfn(s.retriever)
	if err != nil {
		return 0, err
	}
	i, err := r.Read(p)
	r.Close()
	return i, err
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

var StdoutStore = &StoreMaker{"STDOUT", stdoutStore}

func stdoutStore([]string) Store {
	rs := []string{"default", "STDOUT"}
	return &store{
		ofn: func(c *Container, w io.WriteCloser) ([]string, error) {
			b, err := c.MarshalJSON()
			if err != nil {
				return nil, err
			}
			_, err = w.Write(b)
			w.Close()
			return rs, err
		},
		wfn: func(c *Container) (io.WriteCloser, error) {
			return os.Stdout, nil
		},
		ifn: func(r string, n int64, rr io.ReadCloser) (*Container, error) {
			return nil, FunctionNotImplemented("In Function", "STDOUT")
		},
		rfn: func(rt *retriever) (io.ReadCloser, int64, error) {
			return nil, 0, FunctionNotImplemented("Read Function", "STDOUT")
		},
		retriever: newRetriever(rs...),
	}
}

var YamlStore = &StoreMaker{"yaml", yamlStore}

func yamlStore(rs []string) Store {
	return &store{
		ofn: func(c *Container, w io.WriteCloser) ([]string, error) {
			retrieval := c.ToList("store.retrieval.string")
			y, err := yaml.Marshal(&c)
			if err != nil {
				return nil, err
			}
			_, err = w.Write(y)
			w.Close()
			return retrieval, err
		},
		wfn: fileFromContainer,
		ifn: func(r string, n int64, rr io.ReadCloser) (*Container, error) {
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
		rfn:       readCloserFromRetriever,
		retriever: newRetriever(rs...),
	}
}

var (
	JsonStore  = &StoreMaker{"json", JsonStorer(regular)}
	JsonFStore = &StoreMaker{"jsonf", JsonStorer(indented)}
)

type jsonMarshaler func(*Container) ([]byte, error)

func regular(c *Container) ([]byte, error) {
	return json.Marshal(&c)
}

func indented(c *Container) ([]byte, error) {
	return json.MarshalIndent(&c, "", "    ")
}

func JsonStorer(jm jsonMarshaler) StoreFn {
	return func(rs []string) Store {
		return &store{
			ofn: func(c *Container, w io.WriteCloser) ([]string, error) {
				retrieval := c.ToList("store.retrieval.string")
				j, err := jm(c)
				if err != nil {
					return nil, err
				}
				_, err = w.Write(j)
				w.Close()
				return retrieval, err
			},
			wfn: fileFromContainer,
			ifn: func(r string, n int64, rr io.ReadCloser) (*Container, error) {
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
			rfn:       readCloserFromRetriever,
			retriever: newRetriever(rs...),
		}
	}
}

func insufficient(s []string, i int) error {
	if len(s) < i {
		return MalformedRetrievalStringError(s, "expected length equal to or greater than three")
	}
	return nil
}

func fileFromContainer(c *Container) (io.WriteCloser, error) {
	rs := c.ToList("store.retrieval.string")
	if err := insufficient(rs, 3); err != nil {
		return nil, err
	}
	ext, loc, fil := rs[0], rs[1], rs[2]
	p := filepath.Join(loc, fmt.Sprintf("%s.%s", fil, ext))
	fl, err := Open(p)
	if err != nil {
		return nil, err
	}
	fl.Truncate(0)
	return fl, nil
}

var ReaderRetrievalError = Drror("unable to find readcloser: %s").Out

func readCloserFromRetriever(rt *retriever) (io.ReadCloser, int64, error) {
	rs := rt.Retrieval()
	if err := insufficient(rs, 3); err != nil {
		return nil, 0, err
	}
	ext, dir, file := rs[0], rs[1], rs[2]
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
