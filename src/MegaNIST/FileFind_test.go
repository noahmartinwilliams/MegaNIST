package main

import "testing"
import "sort"

func TestFindFiles(t *testing.T) {
	dirName := "dir"
	outputc := FindFiles(dirName, false)
	list := make([]string, 3)
	x := 0
	for entry := range outputc {
		list[x] = entry
		x = x + 1
	}

	sort.Strings(list)
	if list[0] != "dir/a.txt" {
		t.Errorf("list[0] is not \"dir/a.txt\". Got \"%s\".", list[0])
	}
	if list[1] != "dir/b.txt" {
		t.Errorf("list[1] is not \"dir/b.txt\". Got \"%s\".", list[1])
	}
	if list[2] != "dir/dir2/c.txt" {
		t.Errorf("list[2] is not \"dir/dir2/c.txt\". Got \"%s\".", list[2])
	}
}
