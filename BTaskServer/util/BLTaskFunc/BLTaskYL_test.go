package BLTaskFunc

import (
	"fmt"
	"testing"
)

func TestGetOrderList(t *testing.T) {
	res := GetOrderListYL(1, 30, 1)
	fmt.Println(res)
}
