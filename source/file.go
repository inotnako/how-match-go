package source

import "io/ioutil"

func NewFileSourcer() Sourcer {
	return &file{}
}

type file struct{}

func (f *file) Get(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
