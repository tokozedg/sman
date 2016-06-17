package sman

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type file struct {
	name, content string
}

var testFiles = []file{

	// Valid snippet yml
	{"single.yml",
		`no_placeholder:
      do: copy
      tags:
      - single
      command: test command`},
	{"examples.yml",
		`single_placeholder:
        do: copy
        tags:
          - tag1
        command: echo <<fname>>
     multiple_placeholders:
        do: copy
        tags:
          - tag2
        command: echo <<fname>> <<lname>>
     ext_command:
        do: copy`},
	// Test command file read
	{"examples/ext_command", "test command"},
	// Junk file
	{"junk", ""},
}

const testPath = "./testdata/"

func makeTestFiles(t *testing.T) {
	err := os.MkdirAll(testPath+"/examples", 0770)
	if err != nil {
		t.Errorf("makeDir: %v", err)
		return
	}
	for _, f := range testFiles {
		err := ioutil.WriteFile(testPath+f.name, []byte(f.content), 0644)
		if err != nil {
			t.Errorf("makeTree: %v", err)
			return
		}
	}
}

func cleanTestFiles(t *testing.T) {
	if err := os.RemoveAll(testPath); err != nil {
		t.Errorf("removeTree: %v", err)
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name string
		p    string
		want string
	}{
		{"abs path", "/home/user/Documents", "/home/user/Documents"},
		{"home", "~/snippets", os.Getenv("HOME") + "/snippets"},
		{"cleaned", "/home//user/Documents/", "/home/user/Documents"},
	}
	for _, tt := range tests {
		if got := expandPath(tt.p); got != tt.want {
			t.Errorf("%q. expandPath() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestYmlFiles(t *testing.T) {
	tests := []struct {
		name      string
		dir       string
		wantFiles []string
	}{
		{"valid files", testPath, []string{"examples", "single"}},
	}
	makeTestFiles(t)
	for _, tt := range tests {
		if gotFiles := ymlFiles(tt.dir); !reflect.DeepEqual(gotFiles, tt.wantFiles) {
			t.Errorf("%q. ymlFiles() = %v, want %v", tt.name, gotFiles, tt.wantFiles)
		}
	}
	defer cleanTestFiles(t)
}

func TestUnmarshalFile(t *testing.T) {
	tests := []struct {
		name            string
		file            string
		wantSnippetsMap map[string]Snippet
	}{
		{"unmarshal single.yml", testPath + "single.yml",
			map[string]Snippet{
				`no_placeholder`: Snippet{
					Command: "test command",
					Tags:    []string{"single"},
					Do:      "copy",
				},
			},
		},
	}
	makeTestFiles(t)
	for _, tt := range tests {
		if gotSnippetsMap := unmarshalFile(tt.file); !reflect.DeepEqual(gotSnippetsMap, tt.wantSnippetsMap) {
			t.Errorf("%q. unmarshalFile() = %v, want %v", tt.name, gotSnippetsMap, tt.wantSnippetsMap)
		}
	}
	cleanTestFiles(t)
}
