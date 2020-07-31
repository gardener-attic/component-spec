// Code generated for package jsonscheme by go-bindata DO NOT EDIT. (@generated)
// sources:
// ../../../../component-descriptor-v2-Schema.yaml
package jsonscheme

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _ComponentDescriptorV2SchemaYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x58\xcd\x8a\xe3\x38\x10\xbe\xfb\x29\x0a\xd2\xa0\x5d\xe8\xb4\x97\xbe\x2c\xe4\xb2\x2c\xdd\xb0\xec\x69\xa1\xd9\x99\xcb\xd0\x07\x45\xae\x24\x6a\x6c\xc9\x23\xc9\x49\x87\x61\xde\x7d\x90\x65\xc9\x72\x1c\x3b\xf6\xa4\x0f\x33\xf4\x9c\x82\x7e\xea\xef\xab\xaf\x4a\xe5\xdc\xf0\x6c\x05\x64\x67\x4c\xa9\x57\x69\xba\xa5\x2a\x43\x81\xea\x8e\xe5\xb2\xca\x52\xcd\x76\x58\x50\x9d\x32\x59\x94\x52\xa0\x30\xcb\x0c\x35\x53\xbc\x34\x52\x2d\xf7\xf7\x24\xb9\x71\x37\x22\x0d\x2f\x5a\x8a\xa5\xdb\xbd\x93\x6a\x9b\x66\x8a\x6e\xcc\xf2\x8f\x3f\x1b\x5d\x0b\x92\x78\x15\x5c\x8a\x15\x90\x7f\x1a\x8b\xf0\xe0\x6d\xc0\x63\xb0\x01\xfb\x7b\x70\x72\x56\x6c\xc3\x05\xb7\x52\x7a\x95\x00\x14\x68\xa8\xfd\x05\x30\xc7\x12\x57\x40\xe4\xfa\x05\x99\x21\xf5\x56\xd7\x44\xf0\x1e\x5a\xef\x6b\xf9\x8c\x1a\xea\x04\x14\x7e\xae\xb8\xc2\xcc\x69\x04\x58\x02\x71\x76\x3f\xa2\xd2\x5c\x0a\x77\xab\x54\xb2\x44\x65\x38\x6a\x7f\xaf\x73\xc9\x6f\x06\x97\xb4\x51\x5c\x6c\x49\x92\x58\x03\xa5\xd4\xdc\x48\x75\x7c\x90\xc2\xe0\xab\x19\xf4\xfd\x8c\x2b\x6b\xaa\xf1\x83\xca\x87\x9c\x68\x8e\xc7\xcc\x53\xc6\x50\xeb\x89\x78\x59\x7d\xf5\x2d\xd8\x48\xd5\x88\xa2\x86\xdf\xec\x0a\x5f\x0d\x0a\x1b\xac\xfe\x7d\xd0\x5d\x2b\x5a\x5b\xdd\x72\xb3\xab\xd6\x7f\x8f\xdb\x1e\x54\x10\x96\x16\xbb\x10\x7e\xb3\xb3\x19\x02\xa3\xb6\x30\x84\x84\xdf\x46\x51\x15\x2b\xf8\x44\x9c\x83\xe4\xb9\x39\x68\x0c\x5d\x10\x57\xb8\x19\x83\x5a\xcb\x4a\x31\x7c\x0c\x6c\x9d\x11\xb8\xa0\x05\x86\x85\x95\x08\x0b\x97\x84\x81\x90\xad\xd8\x05\x9f\xe7\xa2\x12\x20\x89\x89\x53\xaf\xc5\xf1\xbf\x28\x7c\xeb\xdb\x8d\x05\x04\xc8\x22\x8d\x2a\x34\x75\x72\xe4\xf2\xc5\x98\x23\x35\x7e\xa1\x5a\x9f\x70\x83\x0a\x05\xc3\x89\xb4\xa5\x36\x33\x4e\x02\x8c\x04\xda\x6a\x1a\x24\x9a\x45\x2e\xa2\xd5\x7e\xbc\xd8\xc7\x70\x86\x05\x50\x66\x2a\x9a\xe7\xc7\x55\x6b\x78\x59\x57\xd1\x21\x05\x5d\x22\xe3\x34\x07\x85\xf6\x3e\xab\x43\x6f\x34\xed\xc7\x9b\x47\x47\xb1\xc2\x9c\xbe\x62\x06\x1a\x8b\x3d\xaa\xbf\x5c\x63\x71\x84\xfb\x3f\x64\x78\x66\x7d\x7b\x05\x7a\x2e\x4a\xb0\xa8\xe5\x73\xc9\x68\xfe\xe4\x95\xdc\x82\xd9\x71\x0d\x45\xa5\x0d\x14\xd4\xb0\x5d\x94\x05\xed\x63\x1d\xac\xf4\x98\x35\xdf\xc5\xf3\x0b\x60\x4e\xaf\x86\x53\xde\x7b\xee\x9e\xe3\x78\x02\x20\x19\xff\xb7\xa0\x5b\xbc\xba\xd7\x71\xab\x25\x30\xff\x4d\x9a\x9c\x64\xfc\x09\xb7\x5c\x1b\x75\x0c\x65\xdd\x35\x33\xd6\xce\x7c\x64\x3e\xc3\x73\x62\x1b\x2b\xaf\x1f\x3b\xfd\x1d\xf4\xea\xf8\x07\x3b\xe2\xd9\xae\xd6\xe5\x43\x0d\xe4\x01\xaf\x7f\x09\xab\xe1\x21\x60\x56\x50\x07\x6c\x1f\xbd\x6a\x7c\x72\x38\xe0\xfa\xbd\xa6\x3e\x46\x69\x52\xd6\x43\x8a\xdd\xf4\x63\xc7\x5a\xce\xae\x49\xfa\xb5\xe3\x8d\xf3\x80\x3c\x47\xee\xbc\xd7\x64\xb6\x58\xcc\x48\x68\x27\x85\xdd\xe9\x64\xf2\x50\x72\xdd\x14\x12\x06\xe0\xce\xc7\x83\x8e\x0e\x4b\x25\xf7\x3c\x43\x15\x6d\x75\x1e\xf4\x7a\xa7\x3f\x54\xc5\xa7\xdd\x27\x3c\x3a\xb0\xd3\xbe\x12\xbd\xb3\x9f\x61\x34\x72\x02\x7d\xe4\x7a\xca\xa8\x52\xf4\xd8\xf2\x85\x1b\x2c\x74\x3c\xe0\x9e\x25\x46\x4f\xaf\x57\xe0\xb3\x31\x99\x97\x5c\x38\x88\xc9\x6d\x0b\x77\x20\x69\x83\xfa\xdb\xf8\x7c\xfa\x61\xe2\xe5\xcf\x70\xe3\x6d\x0c\xf6\x15\x7b\x0d\x5d\xc2\xcd\xb6\xd6\xfb\x08\x19\xfc\xba\x88\xa7\x63\x32\x45\xe0\x74\xdc\x99\x24\x14\xbd\x91\x93\xee\x9f\xb4\x62\x2f\xd3\xab\xb6\x5f\xc0\xd4\x42\xe7\x08\x0a\x5f\xbe\x26\x49\x72\xd2\x7e\xe3\xde\xba\x04\x52\xa0\xfb\x7b\x27\xee\x7f\x24\xe9\x36\xaf\xf6\x6f\xa4\xb3\x0e\x79\x15\x27\x6d\x7f\x9c\xee\x24\xf9\x16\x00\x00\xff\xff\xe4\x05\x28\xf2\x55\x13\x00\x00")

func ComponentDescriptorV2SchemaYamlBytes() ([]byte, error) {
	return bindataRead(
		_ComponentDescriptorV2SchemaYaml,
		"../../../../component-descriptor-v2-Schema.yaml",
	)
}

func ComponentDescriptorV2SchemaYaml() (*asset, error) {
	bytes, err := ComponentDescriptorV2SchemaYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../../../component-descriptor-v2-Schema.yaml", size: 4949, mode: os.FileMode(420), modTime: time.Unix(1596715729, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"../../../../component-descriptor-v2-Schema.yaml": ComponentDescriptorV2SchemaYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"..": &bintree{nil, map[string]*bintree{
		"..": &bintree{nil, map[string]*bintree{
			"..": &bintree{nil, map[string]*bintree{
				"..": &bintree{nil, map[string]*bintree{
					"component-descriptor-v2-Schema.yaml": &bintree{ComponentDescriptorV2SchemaYaml, map[string]*bintree{}},
				}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
