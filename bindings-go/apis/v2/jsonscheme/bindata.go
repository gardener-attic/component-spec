// Code generated for package jsonscheme by go-bindata DO NOT EDIT. (@generated)
// sources:
// ../../../../component-descriptor-v2-schema.yaml
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

var _ComponentDescriptorV2SchemaYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x58\x4f\x8b\xe3\x36\x14\xbf\xfb\x53\x3c\xc8\x80\x5a\x58\xaf\xcb\x5e\x0a\xb9\x94\xb2\x0b\xa5\x50\xb6\x30\xb4\xbd\x94\x3d\x28\xf2\x4b\xa2\xc5\x96\x5c\x49\x4e\xc7\x94\x7e\xf7\x22\xc9\xb2\xa5\x38\xf6\x38\x9b\x74\x77\x4a\x4f\x83\x3c\xef\xf7\x7b\xff\x9f\x9e\xf2\xc0\xcb\x2d\x90\xa3\x31\x8d\xde\x16\xc5\x81\xaa\x12\x05\xaa\xd7\xac\x92\x6d\x59\x68\x76\xc4\x9a\xea\x82\xc9\xba\x91\x02\x85\xc9\x4b\xd4\x4c\xf1\xc6\x48\x95\x9f\xde\x90\xec\xc1\x4b\x44\x0c\x1f\xb5\x14\xb9\xff\xfa\x5a\xaa\x43\x51\x2a\xba\x37\xf9\x37\xdf\xf6\x5c\x1b\x92\x05\x0a\x2e\xc5\x16\xc8\x0f\xbd\x46\x78\x1b\x74\xc0\xbb\x41\x07\x9c\xde\x80\xc7\x59\xd8\x9e\x0b\x6e\x51\x7a\x9b\x01\xd4\x68\xa8\xfd\x0b\x60\xba\x06\xb7\x40\xe4\xee\x23\x32\x43\xdc\xa7\x54\xc5\x60\x3d\x8c\xd6\x3b\x7c\x49\x0d\xf5\x00\x85\x7f\xb4\x5c\x61\xe9\x19\x01\x72\x20\x5e\xef\x6f\xa8\x34\x97\xc2\x4b\x35\x4a\x36\xa8\x0c\x47\x1d\xe4\x12\xa1\xf0\x71\x30\x49\x1b\xc5\xc5\x81\x64\x19\x40\x45\x77\x58\xcd\xda\x7b\x41\xbd\xa0\x35\x92\xf1\x78\xa2\x55\x8b\x8e\x49\x61\x23\x35\x37\x52\x75\x6f\xa5\x30\xf8\x64\xae\x61\xdd\x51\x8d\xbf\xaa\x2a\x22\xb6\xb8\x39\xef\x7a\xe9\x59\xbf\xa2\x8f\xcb\x22\x00\x28\xda\x7a\x0b\xbf\x13\xc9\xf8\x23\x1e\xb8\x36\xaa\x23\x1f\xac\x3b\x94\x31\xd4\x7a\x65\x26\xad\x41\x4e\x0a\xf6\x52\xf5\x50\xd4\xf0\x95\x3d\xe1\x93\x41\x61\xd3\xa0\xbf\x9e\x75\xdf\x3b\x9b\x01\x1c\xb8\x39\xb6\xbb\xef\x97\x75\xcf\x12\x0c\x47\x9b\x8b\x34\x9c\x0a\xf7\x73\xd1\xbc\x2a\x4e\xde\x40\xf2\xa1\xff\x47\xaf\xe8\x19\xb8\xc2\xfd\x52\x0d\x0a\x29\xf0\x16\x97\x6f\x74\xe9\xbd\x14\xe8\x73\xae\x65\xab\x18\xbe\x1b\x1a\xfa\x0a\x73\x6c\x5b\x0c\x07\x8b\x18\x0e\xbe\x1a\x66\x0c\xb5\xb0\x7b\x96\xf1\x81\x9b\x21\x37\xae\xb5\xf5\x04\x4a\x95\xa2\xdd\x88\xe4\x06\xeb\x48\x08\xe0\xc1\x66\x0b\xc8\xa6\x88\x06\x5b\xe1\xb8\x02\x28\x6e\x0d\x77\x16\xdd\xcf\xfb\x98\x22\xbf\x4c\xe2\x71\xe4\x79\xc1\xb8\x0b\x56\x88\xdb\x11\x1f\x84\x33\x80\x61\xac\x3e\xe2\x1e\x15\x0a\x86\x2b\xbb\x98\xda\x42\xf5\x08\x30\x12\xe8\xc8\xb4\x76\x1a\x0e\x80\xf7\x67\x53\x72\x79\x5a\x2f\x55\x01\x6c\x80\x32\xd3\xd2\xaa\xea\xb6\xa3\x41\xb9\x1b\x36\x7f\x16\xa0\x1b\x64\x9c\x56\xa0\xd0\xca\x33\x17\x90\x9e\xe9\xb4\x3c\xfd\x13\x62\x85\x15\x7d\xc2\x12\x34\xd6\x27\x54\xdf\xfd\x7b\x15\xe4\x2e\x0a\xdf\x67\xbf\x0c\x85\x7d\xe5\x7c\x0d\x04\x7a\xf5\x25\xd5\x87\x1f\x36\x0e\x5f\x49\xe6\x23\xe6\x59\x5e\x81\x39\x72\x0d\x75\xab\x0d\xd4\xd4\xb0\x63\x94\x77\x1d\xa2\xb8\x30\x6a\x2b\x6a\x86\xdc\xba\x4f\x71\x9d\x7f\x52\xcb\x3f\x93\xb9\xf5\x83\x21\x18\xb7\x7a\x7e\xb8\xd0\x90\x57\x40\xec\xad\xa5\x04\xad\x3e\xff\x34\x09\x88\x4b\x93\x23\x03\x90\x8c\xff\x58\xd3\xc3\x4d\x17\x86\x3b\x72\xcb\x32\x8c\x88\xbb\xdc\x24\xe9\x12\xd1\x47\x24\x51\xb3\x74\x0d\x4a\xc6\x7f\xa2\x1d\xaa\xab\x3c\x4b\xfc\xca\x81\x54\x96\xe1\x1e\xce\x00\x09\xf6\x10\x18\x6b\xa0\x43\x35\x21\x10\x6d\xbd\xb3\x3a\xa3\xe4\x3c\xf6\xbd\x75\xcb\x62\x19\x8f\xcc\x4b\xf9\x7b\x39\x5d\x96\x14\x80\xf3\xff\xcb\x75\x4d\x82\x48\x7b\xc5\x65\x68\xbc\x2b\x6f\x68\x9d\x36\xac\x96\x37\xf6\x8b\x35\x66\x08\x55\xbb\xb0\x46\xda\xcd\xd8\x3e\xc6\x38\xfb\x82\x7b\x62\x6f\x81\x5f\x15\xfb\xc3\xff\xb5\xd0\xc7\x58\xbc\x84\x3a\x4f\x6a\x23\xdd\x01\x57\xaf\x7e\x57\xef\x7a\xd3\xbc\x4d\x5e\xc0\x3a\xfa\x67\xa3\xe4\x89\x97\x61\x36\xfb\x97\x7c\xbc\xc5\xa4\xfb\xe3\x70\x61\xe8\x84\x3f\x41\xfc\x77\xf6\xc8\x69\x60\xee\x53\x29\x13\xde\x40\x10\x82\xbd\xba\x9e\xb9\xe8\xd7\x9d\xcf\xba\xfa\xf4\xe9\xbc\x0f\xf3\xf9\xe3\x35\xe0\x2f\xd4\xd4\x7d\x14\x4e\x89\xc7\xb5\xf3\x53\x3d\x9b\x3c\x25\x67\x1f\x7d\xf1\x2b\x82\xac\x01\x9c\xef\x26\xab\x40\x67\x63\xde\xcd\x96\xcb\x21\x85\xbf\xfe\xce\xb2\xec\x6c\xd0\xc4\x53\x24\x07\x52\xa3\xff\x5d\x2f\xee\x74\x92\xa5\x7d\x3c\xfe\x7e\x78\xd1\xa0\x40\x71\x36\xe0\x96\x13\x44\xb2\x7f\x02\x00\x00\xff\xff\x6c\xfe\x61\x9d\x4e\x15\x00\x00")

func ComponentDescriptorV2SchemaYamlBytes() ([]byte, error) {
	return bindataRead(
		_ComponentDescriptorV2SchemaYaml,
		"../../../../component-descriptor-v2-schema.yaml",
	)
}

func ComponentDescriptorV2SchemaYaml() (*asset, error) {
	bytes, err := ComponentDescriptorV2SchemaYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../../../component-descriptor-v2-schema.yaml", size: 5454, mode: os.FileMode(420), modTime: time.Unix(1603809437, 0)}
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
	"../../../../component-descriptor-v2-schema.yaml": ComponentDescriptorV2SchemaYaml,
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
					"component-descriptor-v2-schema.yaml": &bintree{ComponentDescriptorV2SchemaYaml, map[string]*bintree{}},
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
