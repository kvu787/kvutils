package util

import (
	"encoding/json"
	"testing"
)

func TestRealDeal(t *testing.T) {
	_, err := ConvertFilesToNode(`C:\Users\kevin\wksp`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConvertFilesToNode(t *testing.T) {
	testdirNode, err := ConvertFilesToNode(`testdata\testdir`)
	if err != nil {
		t.Fatal(err)
	}
	if *testdirNode.Name != "testdir" {
		t.Fatal("wrong name")
	}
	if len(testdirNode.Children) != 3 {
		t.Fatal("wrong # of children")
	}

	if testdirNode.String() != `testdir/
  c/
    f/
      g - hello world g
    d - hello world d
    e - hello world e
  a - hello world a
  b - hello world b` {
		t.Fatal("wire representation -> struct failed")
	}
}

func TestJsonEncodeNode(t *testing.T) {
	testdirNode, err := ConvertFilesToNode(`testdata\testdir`)
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, err := json.Marshal(testdirNode)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: fix this so that JSON key reordering doesn't fail the test
	if string(jsonBytes) != `{"Name":"testdir","Data":null,"Children":[{"Name":"c","Data":null,"Children":[{"Name":"f","Data":null,"Children":[{"Name":"g","Data":"aGVsbG8gd29ybGQgZw==","Children":null}]},{"Name":"d","Data":"aGVsbG8gd29ybGQgZA==","Children":null},{"Name":"e","Data":"aGVsbG8gd29ybGQgZQ==","Children":null}]},{"Name":"a","Data":"aGVsbG8gd29ybGQgYQ==","Children":null},{"Name":"b","Data":"aGVsbG8gd29ybGQgYg==","Children":null}]}` {
		t.Fatal("JSON is wrong")
	}
	var testdirNode2 Node
	err = json.Unmarshal(jsonBytes, &testdirNode2)
	if err != nil {
		t.Fatal(err)
	}

	if testdirNode2.String() != `testdir/
  c/
    f/
      g - hello world g
    d - hello world d
    e - hello world e
  a - hello world a
  b - hello world b` {
		t.Fatal("wire representation -> struct failed")
	}
}
