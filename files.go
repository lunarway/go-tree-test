package treetest

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type File struct {
	fileInfo os.FileInfo
	path     string
	parent   *directory
}

func (file File) GetPath() string {
	return file.path
}

type directory struct {
	fileInfo    os.FileInfo
	path        string
	parent      *directory
	files       map[string]*File
	directories map[string]*directory
	test        SpecTest
}

func getFiles(path string) (*directory, error) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return &directory{}, err
	}
	dir := &directory{
		fileInfo:    fileinfo,
		path:        path,
		files:       map[string]*File{},
		directories: map[string]*directory{},
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for i := range files {
		f := files[i]
		if f.IsDir() {
			sub, err := getFiles(filepath.Join(path, f.Name()))
			if err != nil {
				return nil, err
			}
			dir.directories[f.Name()] = sub
		} else {
			dir.files[f.Name()] = &File{
				parent:   dir,
				fileInfo: f,
				path:     filepath.Join(path, f.Name()),
			}
		}
	}
	for n := range dir.directories {
		d := dir.directories[n]
		d.parent = dir
		dir.directories[n] = d
	}
	return dir, nil
}

func (dir *directory) findUpwards(fileName string) *File {
	if file, ok := dir.files[fileName]; ok {
		return file
	}
	if dir.parent != nil {
		return dir.parent.findUpwards(fileName)
	}
	return nil
}

func (dir *directory) findAllFilesUpwards(fileNames ...string) (bool, map[string]File) {
	foundAll := true
	files := make(map[string]File)
	for i := range fileNames {
		fileName := fileNames[i]
		foundFile := dir.findUpwards(fileName)
		if foundFile == nil {
			foundAll = false
			continue
		}
		files[fileName] = *foundFile
	}
	return foundAll, files
}

func (dir *directory) findValidDirectoriesWith(fileNames ...string) []*directory {
	return dir.findDirectoriesWith(append(fileNames, "txcore.json", "user.json", "env.json")...)
}

func (dir *directory) getAllDirectories() []*directory {
	dirs := []*directory{dir}
	for _, sub := range dir.directories {
		dirs = append(dirs, sub.getAllDirectories()...)
	}
	return dirs
}

func (dir *directory) findDirectoriesWith(fileNames ...string) []*directory {
	var dirs []*directory
	for _, sub := range dir.directories {
		dirs = append(dirs, sub.findDirectoriesWith(fileNames...)...)
	}

	for _, fileName := range fileNames {
		foundFile := dir.findUpwards(fileName)
		if foundFile == nil {
			return dirs
		}
	}

	return append(dirs, dir)
}

func (dir *directory) mustGetFile(fileName string) *File {
	if f, ok := dir.files[fileName]; ok {
		return f
	}
	panic(fmt.Sprintf("Could not find File %s in %s", fileName, dir.path))
}

func (dir *directory) getTestSpec() SpecTest {
	if dir.test == nil {
		dir.test = DefineTest(dir.path, *dir)
	}
	return dir.test
}

func (dir *directory) getAllDefinedTestSpecs() []SpecTest {
	var testSpecs []SpecTest
	sort.Slice(testSpecs, func(i, j int) bool {
		return strings.Compare(testSpecs[i].GetName(), testSpecs[j].GetName()) < 0
	})
	for _, sub := range dir.directories {
		testSpecs = append(testSpecs, sub.getAllDefinedTestSpecs()...)
	}
	if dir.test != nil {
		testSpecs = append(testSpecs, dir.test)
	}
	return testSpecs
}

func (f *File) ReadData(withObject interface{}) (interface{}, error) {
	bytes, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, errors.Errorf("failed reading File: %s", f.path)
	}
	err = json.Unmarshal(bytes, withObject)
	if err != nil {
		return nil, errors.Errorf("failed unmarshalling %T for '%s'. Error:\n%s", withObject, f.path, err)
	}
	return withObject, nil
}
func (f *File) MustReadData(withObject interface{}) interface{} {
	_, err := f.ReadData(withObject)
	if err != nil {
		panic(err)
	}
	return withObject
}
func (f *File) ReadString() (string, error) {
	bytes, err := ioutil.ReadFile(f.path)
	if err != nil {
		return "", errors.Errorf("failed reading File: %s", f.path)
	}
	return string(bytes), nil
}

func (f *File) MustReadString() string {
	str, err := f.ReadString()
	if err != nil {
		panic(err)
	}
	return str
}

type Files map[string]File

func (files Files) MustReadString(fileName string) string {
	file, ok := files[fileName]
	if !ok {
		var fileNames []string
		for fn := range files {
			fileNames = append(fileNames, fn)
		}
		panic(fmt.Sprintf("FileName '%s' is not defined! Available: %s", fileName, strings.Join(fileNames, ", ")))
	}
	return file.MustReadString()
}

func (files Files) MustReadData(fileName string, withObject interface{}) interface{} {
	file, ok := files[fileName]
	if !ok {
		var fileNames []string
		for fn := range files {
			fileNames = append(fileNames, fn)
		}
		panic(fmt.Sprintf("FileName '%s' is not defined! Available: %s", fileName, strings.Join(fileNames, ", ")))
	}
	return file.MustReadData(withObject)
}
