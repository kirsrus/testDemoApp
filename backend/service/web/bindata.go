/**
!!!! Это специальный переходной модуль для go-bindata. НЕ ИЗМЕНЯТЬ (кроме импорта bindata)
Нужен для корректого отображения рабочих функций, т.к. из-за обольшого объёма
конечного файла bindaoa.go IDE отказывается его индексировать.
*/
package web

import (
	"TestDemoApp/bindata"
	"net/http"
	"os"
)

// AssetFile return a http.FileSystem instance that data backend by asset
func AssetFile() http.FileSystem {
	return bindata.AssetFile()
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	return bindata.Asset(name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	return bindata.MustAsset(name)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	return bindata.AssetInfo(name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	return bindata.AssetNames()
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
// AssetDir("foo.txt") and AssetDir("nonexistent") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	return bindata.AssetDir(name)
}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	return bindata.RestoreAsset(dir, name)
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	return bindata.RestoreAssets(dir, name)
}
