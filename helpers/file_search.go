package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
)

type FileSearch interface {
	Exists(string) bool
	FileSize(string) int
	Get(string) ([]byte, error)
	Save(string, []byte) error
	Size() int64
	Count() int
}

type fsHelper struct {
	dir     string
	changed bool
	files   map[string]int64
	size    int64
}

func NewFiles(dir string) FileSearch {
	f := &fsHelper{
		dir:     dir,
		changed: true,
		files:   make(map[string]int64),
	}
	f.update()
	return f
}

func (h *fsHelper) Exists(fName string) bool {
	h.update()
	_, ok := h.files[fName]
	return ok
}

func (h *fsHelper) FileSize(fName string) int {
	h.update()
	return h.files[fName]
}

func (h *fsHelper) Get(fName string) ([]byte, error) {
	h.update()
	if _, ok := h.files[fName]; !ok {
		return nil, fmt.Errorf("unable to find  file: %s", fName)
	}

	buf, err := ioutil.ReadFile(fmt.Sprintf("./%s/%s", h.dir, fName))
	if err != nil {
		return nil, err
	}
	return buf, err
}

func (h *fsHelper) Save(fName string, data []byte) error {
	err := ioutil.WriteFile(fmt.Sprintf("./%s/%s", h.dir, fName), data, os.ModePerm)
	if err != nil {
		return err
	}

	h.changed = true
	return nil
}

func (h *fsHelper) Size() int64 {
	h.update()
	return h.size
}

func (h *fsHelper) Count() int {
	h.update()
	return len(h.files)
}

func (h *fsHelper) update() {
	if !h.changed {
		return
	}

	files, err := ioutil.ReadDir(fmt.Sprintf("./%s", h.dir))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	h.size = 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		h.files[file.Name()] = file.Size()
		h.size += h.files[file.Name()]
	}

	h.changed = false
}
