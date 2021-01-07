package httplib

import (
	"net/http"
	"strconv"
)

func GetLimitAndOffsetFromRequest(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.FormValue("page"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	if offset > 0 && page == 0 {
		page = offset / limit
	}

	return GetLimitAndOffsetFromPageNumber(page, limit)
}

func GetPagingInt32FromRequest(r *http.Request) (int32, int32) {
	page, _ := strconv.Atoi(r.FormValue("page"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	if offset > 0 && page == 0 {
		page = offset / limit
	}

	return GetPagingInt32FromPageNumber(page, limit)
}

func GetPagingInt64FromRequest(r *http.Request) (int64, int64) {
	page, _ := strconv.Atoi(r.FormValue("page"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	if offset > 0 && page == 0 {
		page = offset / limit
	}

	return GetPagingInt64FromPageNumber(page, limit)
}

func GetPagingInt32FromPageNumber(page, limit int) (int32, int32) {
	l, o := GetLimitAndOffsetFromPageNumber(page, limit)
	return int32(l), int32(o)
}

func GetPagingInt64FromPageNumber(page, limit int) (int64, int64) {
	l, o := GetLimitAndOffsetFromPageNumber(page, limit)
	return int64(l), int64(o)
}

func GetLimitAndOffsetFromPageNumber(page, limit int) (int, int) {
	if page == 0 {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	return limit, offset
}
