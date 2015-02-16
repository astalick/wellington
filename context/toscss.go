package context

// #include <stdlib.h>
// #include "sass2scss.h"
import "C"

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"unsafe"
)

// ToScss converts Sass to Scss with libsass sass2scss.h
func ToScss(r io.Reader) string {
	bs, _ := ioutil.ReadAll(r)
	in := C.CString(string(bs))

	defer C.free(unsafe.Pointer(in))

	chars := C.sass2scss(
		// FIXME: readers would be much more efficient
		in,
		// SASS2SCSS_PRETTIFY_1 Egyptian brackets
		C.int(1),
	)

	return C.GoString(chars)
}

func testToScss(t *testing.T) {
	file, err := os.Open("../test/whitespace/one.sass")
	if err != nil {
		t.Fatal(err)
	}
	e := `$font-stack:    Helvetica, sans-serif;
$primary-color: #333;

body {
  font: 100% $font-stack;
  color: $primary-color; }
`
	if s := ToScss(file); s != e {
		t.Errorf("got:\n%s\nwanted:\n%s", s, e)
	}
}
