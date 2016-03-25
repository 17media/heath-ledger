package models

/*
.__                   __  .__              .__             .___
|  |__   ____ _____ _/  |_|  |__           |  |   ____   __| _/ ____   ___________
|  |  \_/ __ \\__  \\   __\  |  \   ______ |  | _/ __ \ / __ | / ___\_/ __ \_  __ \
|   Y  \  ___/ / __ \|  | |   Y  \ /_____/ |  |_\  ___// /_/ |/ /_/  >  ___/|  | \/
|___|  /\___  >____  /__| |___|  /         |____/\___  >____ |\___  / \___  >__|
     \/     \/     \/          \/                    \/     \/_____/      \/
*/

import (
        "fmt"
        "github.com/maxwellhealth/bongo"
        url "net/url"
        "strconv"
)

// ResultSet is an abstract layer of bongo.ResultSet
type ResultSet struct {
        *bongo.ResultSet
}

// PaginationInfo is the metadata for pagination
type PaginationInfo struct {
        Limit    int         `json:"limit"`
        Offset   int         `json:"offset"`
        Next     interface{} `json:"next"`
        Previous interface{} `json:"previous"`
        Total    int         `json:"total"`
}

// DefaultPagination for default limit and offset
func DefaultPagination() *PaginationInfo {
        return &PaginationInfo{Limit: 50, Offset: 0}
}

// Paginate is the actual function to paginate a query
func (r *ResultSet) Paginate(info *PaginationInfo, host string, urlObj *url.URL) (PaginationInfo, error) {
        // Get count on a different session to avoid blocking
        sess := r.Collection.Connection.Session.Copy()

        total, err := sess.DB(r.Collection.Connection.Config.Database).C(r.Collection.Name).Find(r.Params).Count()
        sess.Close()

        if err != nil {
                return *info, err
        }

        r.Query.Skip(info.Offset).Limit(info.Limit)

        info.Total = total
        nextURL := ""
        prevURL := ""

        nextQuery := urlObj.Query()
        prevQuery := urlObj.Query()
        rawURL := urlObj.String()
        nextURLObj, _ := url.Parse(rawURL)
        prevURLObj, _ := url.Parse(rawURL)
        if total > info.Limit+info.Offset {
                nextQuery.Set("offset", strconv.Itoa(info.Offset+info.Limit))
                nextURLObj.RawQuery = nextQuery.Encode()
                nextURL = fmt.Sprintf("%s%s%s", nextURLObj.Scheme, host, nextURLObj.String())
        }

        if info.Offset > 0 {
                newOffset := info.Offset - info.Limit
                if newOffset < 0 {
                        newOffset = 0
                }
                prevQuery.Set("offset", strconv.Itoa(newOffset))
                prevURLObj.RawQuery = prevQuery.Encode()
                prevURL = fmt.Sprintf("%s%s%s", prevURLObj.Scheme, host, prevURLObj.String())
        }

        if nextURL != "" {
                info.Next = nextURL
        } else {
                info.Next = nil
        }

        if prevURL != "" {
                info.Previous = prevURL
        } else {
                info.Previous = nil
        }

        return *info, nil
}
