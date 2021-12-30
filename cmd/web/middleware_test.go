package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)

	switch sw := h.(type) {
	case http.Handler:
		fmt.Printf("the type is %T\n", sw)
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T\n", sw))
	}

}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	switch sw := h.(type) {
	case http.Handler:
		fmt.Printf("the type is %T\n", sw)
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T\n", sw))

	}
}

func TestWriteToConsole(t *testing.T) {
	var myH myHandler

	h := WriteToConsole(&myH)

	switch sw := h.(type) {
	case http.Handler:
		fmt.Printf("the type is %T\n", sw)
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T\n", sw))

	}
}
