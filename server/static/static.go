package static

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/static/css/normalize.css": {
		local:   "server/static/css/normalize.css",
		size:    7797,
		modtime: 1471461643,
		compressed: `
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
`,
	},

	"/static/css/skeleton.css": {
		local:   "server/static/css/skeleton.css",
		size:    11458,
		modtime: 1471461643,
		compressed: `
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
AAD//7viGWrCLAAA
`,
	},

	"/static/js/list.min.js": {
		local:   "server/static/js/list.min.js",
		size:    15785,
		modtime: 1472414777,
		compressed: `
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
FG1iGP5vAAAA//84G4UUqT0AAA==
`,
	},

	"/static/js/table.js": {
		local:   "server/static/js/table.js",
		size:    1432,
		modtime: 1473627215,
		compressed: `
H4sIAAAJbogA/3xUXW8aOxB951dMnIc1CjHk5ublcjdXVxEqrUgjBaRWiiJk1gO4GHtje/lQxX+v7AV2
E2hfEGufM3PmeGamhc68NBosvg2k86jR0ib8bACsuAWLmYMUvgyfvrKcW4fUz6VjFl1utMMRbnyzu8d6
PlEIKQiTFUvUnr0VaLdDVJh5Y2lymaH11xGVVKSJEVtIS/IHBrlU0nkSseu5VAg0wpk1a8cU6pmfwz10
SrlQxmICFXp8NmvaicxdA0A6V6AI9bFMIbc03oTi2NTYHs/m9GjEki9wFNRQVC2QLeDWHjJES8w66I25
pHZofch1fRNDQrjeHz+gUrTTZFJrtP3R4wBSQMUWuB1LcRZ8cwLOLHKPYsz9WcJfJwTc5NKiO4u+PUHn
VupM5lydJ/x9QrC4Mgvcq2+34bMWuIEFbqUArgVUAdkxYIZKuZfOK8sUdy48AqTwQiKJvHY/4G4/4qqQ
dXBN5V0znsop0ItK4uHJ6rHvXt/Vk/wrdV54cH6rMCVLbmdS/9PpEvDbHFOSzTFbTMyGwIqrAlOSXFUP
CFeQENB8GYBo/VgKAlLUPtr3Sal3V+tOnueoxcNcKkGtWZcdGn9rPWox+hq6dNdoHDtTGS4e0HpHuVK1
GYUUNK7h++Og733+jG8FOk+PI4Za5Ebq4GbS5mIpdTtodOyHMzpKDNYdQ0JFuEoh+Y8rlXpbYHIYJstM
jpomn3qjpHUElxPFuBC9FWp/3CVJUJ206vtlD3WoBW2GCoPKYf/p2/j/QXiXKVcOu7XCvZnNFPZia4va
dpoU3hv9+5VDLkvmday3XCS1PBeH/+G8Mvdw2jxYU1H22g5Glfnf9RQZzs0a9lJJtAxQOfwDoy8Fvmc0
do1fAQAA//8XtlT2mAUAAA==
`,
	},

	"/": {
		isDir: true,
		local: "server",
	},

	"/static": {
		isDir: true,
		local: "server/static",
	},

	"/static/css": {
		isDir: true,
		local: "server/static/css",
	},

	"/static/js": {
		isDir: true,
		local: "server/static/js",
	},
}
