package controllers

import (
	"fmt"
	"testing"
)

func TestNewApi(t *testing.T) {

	var pageSize uint64 = 10

	//总记录数
	var total uint64 = 3110 + 1

	//总页数
	var totalPage = total / pageSize

	//如果不能除尽，那么就需要加1
	if total%pageSize != 0 {
		totalPage++
	}
	//当前页（默认为第一页）
	var currentPage uint64 = 1

	var start uint64 = 0

	if total >= (currentPage-1)*pageSize+1 {
		start = total - (currentPage-1)*pageSize - 1
	}
	var end uint64 = 0

	if start >= pageSize {
		end = start - pageSize + 1
	} else {
		end = 0
	}

	fmt.Println("start=", start, "end=", end, "totalPage=", totalPage)
}
