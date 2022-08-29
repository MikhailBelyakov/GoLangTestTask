package common

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

const defaultOrderField = "id"
const defaultSortMethod = "DESC"

type SortStruct struct {
	Order string `json:"order"`
	Sort  string `json:"sort"`
}

func sortMethods() map[string]bool {
	return map[string]bool{
		"asc":  true,
		"ASC":  true,
		"desc": true,
		"DESC": true,
	}
}

func (sortStruct *SortStruct) BindParam(context *gin.Context) {
	orderString := context.Query("order")
	orderSplit := strings.Split(orderString, "_")

	if len(orderSplit) > 1 {
		sortStruct.Order = orderSplit[0]

		if _, inMap := sortMethods()[orderSplit[1]]; !inMap {
			sortStruct.Sort = defaultSortMethod
		}

		sortStruct.Sort = orderSplit[1]
	} else {
		sortStruct.Order = defaultOrderField
		sortStruct.Sort = defaultSortMethod
	}
}

type PaginateStruct struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (paginateStruct *PaginateStruct) BindParam(context *gin.Context) {
	offsetString := context.Query("offset")
	limitString := context.Query("limit")

	offset, err := strconv.Atoi(offsetString)

	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(limitString)

	if err != nil || limit == 0 {
		limit = 20
	}

	paginateStruct.Limit = limit
	paginateStruct.Offset = offset
}
