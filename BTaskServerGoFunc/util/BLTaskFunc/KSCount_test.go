package BLTaskFunc

import (
	"fmt"
	"testing"
)

func TestGetKSCount(t *testing.T) {
	shortUrl := "https://v.kuaishou.com/v0rFDU"
	resMap := GetKSCount(shortUrl, 12)
	fmt.Println(resMap)
}
