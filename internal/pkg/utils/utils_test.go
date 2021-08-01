package utils_test

import (
	"reflect"
	"testing"

	"dns-exporter/internal/pkg/utils"

	"github.com/spf13/afero"
)

func TestWriteToFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	suite := map[string]string{
		"domain1.com":  "zonefile content",
		"domain2.com.": "zonefile content",
	}

	for d, c := range suite {
		filename, err := utils.WriteToFile(d, c, "./", fs)

		if err != nil {
			t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
		}

		result, err := afero.ReadFile(fs, filename)
		if err != nil {
			t.Fatal("error reading exported zonefile:", err)
		}

		if !reflect.DeepEqual(c, string(result)) {
			t.Errorf("\nEXPECTED content: \n%+v\n\nGOT content: \n%+v\n\n", c, string(result))
		}
	}
}

func TestValidateDir(t *testing.T) {
	fs := afero.NewMemMapFs()

	// existing directory
	err := fs.MkdirAll("./exists", 0777)
	if err != nil {
		t.Fatal("error creating directory: ", err)
	}

	r1, err := utils.ValidateDir("./exists", true, fs)
	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}
	if !r1 {
		t.Fatal("\nEXPECTED result: \ntrue\n\nGOT error:", r1)
	}

	// create non existing directory
	r2, err := utils.ValidateDir("./create", true, fs)
	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}
	if r2 {
		t.Fatal("\nEXPECTED result: \nfalse\n\nGOT error:", r2)
	}

	// do not create non existing directory
	r3, err := utils.ValidateDir("./do_not_create", true, fs)
	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}
	if r3 {
		t.Fatal("\nEXPECTED result: \nfalse\n\nGOT error:", r3)
	}
}
