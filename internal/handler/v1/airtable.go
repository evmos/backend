// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/tharsis/dashboard-backend/internal/requester"
	"github.com/valyala/fasthttp"
)

func GetAnnouncements(ctx *fasthttp.RequestCtx) {
	path := "/Announcement?maxRecords=30&sort[0][field]=Start+Date+Time&sort[0][direction]=desc"

	if resp, err := db.RedisGetAirtableRequest(path); err == nil {
		sendResponse(resp, nil, ctx)
		return
	}

	resp, err := requester.MakeAirtableGetRequest(path)
	if err != nil {
		if val, err := db.RedisGetAirtableFallbackRequest(path); err == nil {
			sendResponse(val, nil, ctx)
			return
		}
		sendResponse("unable to get airtable request", err, ctx)
	}

	db.RedisSetAirtableRequest(resp, path)
	db.RedisSetAirtableFallbackRequest(resp, path)

	sendResponse(resp, nil, ctx)
}

func AddAirtableRoutes(r *router.Router) {
	r.GET("/Announcements", GetAnnouncements)
}
