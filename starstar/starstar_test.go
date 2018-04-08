package starstar

import (
	"os"
	"testing"
)

func touch(path string, t *testing.T) error {
	fd, err := os.Create(path)
	fd.Close()
	if err != nil {
		t.Fatalf("%s: could not make as test file", path)
		return err
	}
	return nil
}

func TestExpand(t *testing.T) {
	os.Mkdir("sub", 0777)
	defer os.RemoveAll("sub")
	os.Mkdir("sub/sub", 0777)

	testfiles := []string{
		"sub/sub/foobar.x",
		"sub/barbar.xx",
		"sub/dummy",
	}
	for _, f := range testfiles {
		if touch(f, t) != nil {
			return
		}
	}

	result, err := Expand("./**/*bar*")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	if len(result) < 1 {
		t.Fatal("item not found")
		return
	}
	for _, f := range result {
		println(f)
	}
	if len(result) != 2 {
		t.Fatalf("invalid item count: %d", len(result))
	}

	result, err = Expand("star*")
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, f := range result {
		println(f)
	}
	return
}
