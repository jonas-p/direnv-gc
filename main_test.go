package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestHook(t *testing.T) {
	raw, err := ioutil.ReadFile("shell/hook.sh")
	if err != nil {
		t.Error(err)
	}

	asset, err := Asset("shell/hook.sh")
	if err != nil {
		t.Error(err)
	}

	if res := bytes.Compare(raw, asset); res != 0 {
		t.Errorf("Asset and raw data does not match, did you run go-bindata shell/?")
	}
}
