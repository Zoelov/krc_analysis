package krc

import (
	"fmt"
	"testing"
)

func TestKrc(t *testing.T) {
	filePath := "/Users/mono/Desktop/workspace/zproject/krc_analysis/test/阿美-小草.krc"

	lrc, err := GetLrcFromKrc(filePath)

	if err != nil {
		t.Errorf("got err")
	}
	fmt.Println(lrc)
}
