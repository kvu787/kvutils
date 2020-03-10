package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	ErrIrregularFileFound = errors.New("filetree contains an irregular file")
)

// The wire representation of a filetree:
// UTF-8 encoded JSON.
// File names are Unicode and UTF-8 encoded.
// File contents are byte streams which are base64 encoded and UTF-8 encoded.
// Go's json.Marshal already does this, but you
// can implement an en/decoder in any language.
type Node struct {
	// Go's os package represents filenames as strings,
	// so we assume that all OSes can convert their filenames
	// to Go strings.
	Name     *string
	Data     []byte
	Children []*Node
}

func (node Node) Encode() ([]byte, error) {
	return json.Marshal(node)
}

func (node Node) String() string {
	return node.string(0)
}

func (node Node) string(indent int) string {
	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString(strings.Repeat(" ", indent))
	stringBuilder.WriteString(*node.Name)
	if node.IsDir() {
		stringBuilder.WriteString("/")
		for _, child := range node.Children {
			stringBuilder.WriteString("\n")
			stringBuilder.WriteString(child.string(indent + 2))
		}
	} else {
		stringBuilder.WriteString(" - ")
		stringBuilder.WriteString(string(node.Data))
	}
	return stringBuilder.String()
}

func (node Node) IsDir() bool {
	return node.Data == nil
}

func ConvertFilesToNode(path string) (*Node, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fileInfo.Mode().IsRegular() {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		name := fileInfo.Name()
		return &Node{&name, contents, nil}, nil
	} else if fileInfo.Mode().IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		fileNames, err := file.Readdirnames(-1)
		if err != nil {
			return nil, err
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
		children := []*Node{}
		for _, fileName := range fileNames {
			childPath := filepath.Join(path, fileName)
			node, err := ConvertFilesToNode(childPath)
			if err != nil {
				return nil, err
			}
			children = append(children, node)
		}

		sort.Slice(children, func(i, j int) bool {
			a := children[i]
			b := children[j]
			if a.IsDir() && !b.IsDir() {
				return true
			} else if !a.IsDir() && b.IsDir() {
				return false
			} else {
				result := strings.Compare(*a.Name, *b.Name)
				return result == -1
			}
		})

		name := fileInfo.Name()
		return &Node{&name, nil, children}, nil
	} else {
		return nil, ErrIrregularFileFound
	}
}
