// Copyright 2016 The Linux Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by "esc -private -pkg=schema -include=.*\.json$ ."; DO NOT EDIT.

package schema

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
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
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
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

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/config-schema.json": {
		name:    "config-schema.json",
		local:   "config-schema.json",
		size:    2969,
		modtime: 1625865937,
		compressed: `
H4sIAAAAAAAC/+RWsW7bQAzd9RWCkjGJOnTymnYrkAJG2yEojLNE2Ux1xyuPMioU/vdCJzvx2SfZteEu
XSnyvcdHSuLvJE2zElzBaAXJZJM0e7JgHsmIQgOcPpKpcJFOLRRYYaF81l1XduuKJWjVlSxF7CTPXxyZ
+z76QLzIS1aV3L97n/exm74Oy22Jm+Q5WTDFls35sj47R60WkBeev6+U1kJXS/MXKKSPWSYLLAgum6Rd
O2maFQxKoHwN7JQ6YTQ9nI9XxFpJ96RUAveCGjL/bN2nZKqRJfEIVJjNxRIFCmkYTq1ZKUZl5NR0cqdn
PqyAXT/WUysqUJ34KIliVu2bdyigd/MGwNN0HZBsJhrB35mrj0dm6+NfHHAQGWR+ZfU5H39ZclB+Jha3
X3/LUPk1gMo9dIt8k5dQocFu4V2ulZ165KdeYxzfrIZkhdYN2DfayNbGQ1Lh1hIGG9RP08BT19NzQBDS
jUockXlEaih3T/KoCtPUdYi0i/M9asGjLv/b3r9S3WhwZ/Z/7tZfu6tvxD/QLD7gWe/5JzWH+tqWTCOD
v6YlUyE7xYVR9cmWRD+/TCSVu+TzW2JVzbB0//5bt8kf6z6gBdPog4lntWqBXRZ3PNljzRh+Nsj+mniO
eLCvLtlF3Hq+RCfE7WX/1L3xDA8oegEdd2vsGoqs9+FldHyodxGNs3l7AQZpDQd/vr8AAG2lnfmNGAaZ
E9WgzMB+hm9ZsmE43JvwOHy75sL3Mul2Z538CQAA//9C38scmQsAAA==
`,
	},

	"/content-descriptor.json": {
		name:    "content-descriptor.json",
		local:   "content-descriptor.json",
		size:    1079,
		modtime: 1625865919,
		compressed: `
H4sIAAAAAAAC/5yTsW7cMAyGdz8F4QTIkos6BB2MIEu7d2i3ooNOok5Mz5JK8RBci7x7QcvX2G2RILfZ
xP+Rn2zqVwfQe6yOqQjl1A/QfyqYPuQklhIy6BMmgY9zKDN8LugokLMTca0tLquLOFrFo0gZjHmoOW1a
9Sbzzni2QTbvbk2rXTSO/AmpgzG5YHKnyXXCWtr4P9MbJ8eCSubtAzpptcK5IAth7QfQgwH0I3qyX1q4
lf49r0SEKadNIQfQAmNAxuTQw2LGhF8yBuU8hrp5FrvRE18Yj4ESae9qnqdP7FNr0Vf6+ZqPRoASbI+C
9Y1O/xGhJO9v1xKedljlFQ3HxyJ5x7ZEcuAiuu/1MEJjT1rN5Vp19bVYEeQEV3d2v8tMEsf74U5/rEd/
f3XOd5xdV/4H3tcX7C3sqSlqEALnER4juQgSqc7OMNojbBF8fkz7bD36c+wmk5WbTSnLdDtWim9fdrPs
dIbaEm+G3WzZM/44EKMqff37riz3dL0uHcC37qn7HQAA//9DKIMKNwQAAA==
`,
	},

	"/defs-descriptor.json": {
		name:    "defs-descriptor.json",
		local:   "defs-descriptor.json",
		size:    844,
		modtime: 1625865919,
		compressed: `
H4sIAAAAAAAC/5SST2/TTBDG7/kU826jt0DiOHBAqlWKKnrnUE6t0mi6O7aneP9od6IqVPnuaG03SYtA
cLC1+2jmefwbz9MEQBlKOnIQ9k5VoK6oZsf5liBgFNabDiOIh6+B3BfvBNlRhKuxzUe4DqS5Zo29x3ww
3buoCnIOgLJkGL9tA+0lAMUmp7YiIVVl6QM5/ZyRFj42ZdItWSzZYkOl2aeWB7f5s5cM3ipJZNcc9IAi
FHu8u9vL4gaLH8vibHU4/ncy/b+4Wy9mq6fl/P2Hj7vy78qmqo/YDUnKcENJjuleDVdaAh23QXwTMbSs
Qbekv6eNhaEXfA25yN8/kJY5sOuvIwCcnmPX+MjS2ovqPI/KkLk4/ccJjFyzN5+r29liXaz2ytt3VT5f
FjfL4uzTuljNXhFuYpf+wIfQ8QCRC6GO3sJjy7oFaTmNVGBxC/cExj+6zqMh8+v3Y4y4PcgsZI9zf08K
oGofLea/oDaR1ajvXmCgc17w5XoCqGmkOvcZqtPiIXl3Uh4tcmkxXPdpw3uczCQ/u8nPAAAA///5nDLG
TAMAAA==
`,
	},

	"/defs.json": {
		name:    "defs.json",
		local:   "defs.json",
		size:    1670,
		modtime: 1625865903,
		compressed: `
H4sIAAAAAAAC/7STza6bMBCF9zzFyO2S9oJtbGDb7hMpy6oLSiaJq2AjY6RWEe9e8RNChFuJKneRgGc8
3zmeMbcAgByxKa2qnTKa5EC+4klp1a8aaBs8grtY054vpnXgLgi7GvUXo12hNFo41FiqkyqLoTwceTOA
5NBLABClXTqvAIj7XWOvprTDM9qhckhUSquqrUgOn2KaPsLFrykcUzkEu3Amx2IrmlEpfPA+vsIzuhVP
Yy55ygT3aczJlZDgW4UyShmTNGIiTbiUIooij6Jn15N0+x/T8enQJFlxN8/GBxZJwtbozXPxoTnNeCYk
zdb8zePw8eOUcyE5jySTUZYk1Nf8WOxNz7VLQaNxdyI5fJsCMKeG9EeLfZZ8eFt8cG9Ty+eNXeivvp9G
t9frYvf09t3Ti1c6FPy1DhtnlT5vd3jXGOtf66kq6sOAHf99V8n8+Imle9ykunAOrd5bU6N1CptFEQD5
fIvD7in0ryMEy+fK1G6UfmdTE+tvpoL+1wV/AgAA//96IpqyhgYAAA==
`,
	},

	"/image-index-schema.json": {
		name:    "image-index-schema.json",
		local:   "image-index-schema.json",
		size:    2993,
		modtime: 1625865919,
		compressed: `
H4sIAAAAAAAC/6yWz0/jOhDH7/0rRgGJC5CnJ/QOFeLy9sJpD4v2suJg7EkybGNnx1Ogu+r/vrJN2qRJ
C4Te2rHnO5/vxL/+zAAyg14zNULOZnPIvjZo/3dWFFlkuK1ViXBrDb7AtwY1FaRVnHoeck+9rrBWIa8S
aeZ5/uidvUjRS8dlblgVcvHPVZ5iJymPTJvi53nuGrS6LeljWpqdUyifUyifEmXVYEh1D4+oJcUadg2y
EPpsDsESQJbyvyP7ZCuFh27vKvJQEC4M+GQPPUiFECtDrAxJDJ6SGigPygJZwRI5IkTlCZ7yPuZGqnU5
qFGTpXpZZ3P4dxtTL20shtZpJKuVpQK9+K79Vlkxq1WHXbDuzvuwnbbYl9f2ui30+Fd7HWH8tSTGUOvH
Jhrg0ZC6C2nn3bCn3zsRQyV6yTah+474yMIYyPcHhgskrIU4O3gAV8TFwVggo9VoYGApipwyFiHbYOEv
zKYnl2F3nOQGC7IUKvh8S9JRWA9Nv4czTASy8LAS9JNYRwDJyn9X++Fe+/8ePM2rRlzJqqlIg65Q//TL
GpJCi5sYz4ON8LdRIsgWzq7VonRMUtU38+uwFg2am7Ppfd9dN7u+lrzwb7pSsKCEHqZDwa6G54p0BRLO
leQFarWCBwTjnu3CKYNmOnWk2svcLJQUjush98c280Znh3PvNj60leOYYl2RoJYl404eQOZ6nnp7+PA+
HmoPxye7zw9Cd9rhhcmW2c6E9ZjNY+I5fxyoy6fBLXkMuI3scSALVOE7HLuFW90DmP3Lslt2cG2+2yTA
+k3bT4pJWRm3/EYPZ/v+9Y8MZa2T+KDznz01tgdX3lWdfNZ1RWZjXtpf696zZ9zRpNfZmI3PGAigEXN4
VmZjL8HOE24GcD9bz/4GAAD//yCnv52xCwAA
`,
	},

	"/image-layout-schema.json": {
		name:    "image-layout-schema.json",
		local:   "image-layout-schema.json",
		size:    439,
		modtime: 1625865903,
		compressed: `
H4sIAAAAAAAC/2yPQUvEMBCF7/0VQ/Sg4DYVPOW6pwVhD4IX8VDTaTvLNonJVFik/12SaRXRU5g38+W9
91kBqA6TjRSYvFMG1DGg23vHLTmMcJjaAeGxvfiZ4cmOOLXqLlPXSQYDamQORutT8m4nau3joLvY9rxr
HrRoV8JRtyHJaO0DOruZpYLJtaZsrM/FWEi+BMysfzuhXbUQfcDIhEkZyG2yQyYl8TPGJLVk97fth1yA
74FHhOP+8LvyDbmy8JZ2EgZ6OuNtsS8fbrESR3LDj45unpSBl3UGUPd1UzdqnV/Lu1QAS2kS8X2miN03
8l+PKnNL9RUAAP//k31n5bcBAAA=
`,
	},

	"/image-manifest-schema.json": {
		name:    "image-manifest-schema.json",
		local:   "image-manifest-schema.json",
		size:    921,
		modtime: 1625865903,
		compressed: `
H4sIAAAAAAAC/5ySMW8iMRCF+/0VI0MJ+O501bZXUZxSJEoTpXB2x7uDWNsZmygo4r9HtnHAkCKifTvv
zTdv/dEAiB59x+QCWSNaEHcOzT9rgiKDDOtJDQj/lSGNPsC9w440dSpNL6J97rsRJxWtYwiulXLjrVlm
dWV5kD0rHZa//sqszbKP+mLxrZTWoenKVp9seVpSJJDTkSB7w95hdNuXDXZHzbF1yIHQixbiYQAiRzwi
+3xclq9vfhjJgybc9uDzheghjAhpOZTlkPPgLQeC8qAMkAk4ICeKFH7bZbKG/Uort16tmcjQtJtEC39O
mnovWpIO+YvorNE0nDcwZ9QxNqKhCcvSiOVV/H+ism/VHtmf2wuVYlb7imkdcIqjv099HJVi/ul2gENF
oYyxIb28CuXGus/TFpet9Kj9JdRM9qjJULJU9qawJlLB+Lojxoj19N07rP9JXXED8Nwcms8AAAD//7u3
Dj+ZAwAA
`,
	},

	"/": {
		name:  "/",
		local: `.`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	".": {
		_escData["/config-schema.json"],
		_escData["/content-descriptor.json"],
		_escData["/defs-descriptor.json"],
		_escData["/defs.json"],
		_escData["/image-index-schema.json"],
		_escData["/image-layout-schema.json"],
		_escData["/image-manifest-schema.json"],
	},
}
