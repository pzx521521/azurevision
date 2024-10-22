package main

import (
	"fmt"
	"testing"
)

func TestTrans(t *testing.T) {
	v := NewAzureVision()
	v.Feature = "tags"
	anlyze, _ := v.Anlyze("input.png")
	fmt.Printf("%v\n", anlyze)
}
