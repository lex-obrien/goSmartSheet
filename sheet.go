package goSmartSheet

import (
	"strconv"
	"strings"
	"time"
)

//Sheet represents a Smart Sheet object
type Sheet struct {
	ID                         int64    `json:"id"`
	Name                       string   `json:"name"`
	Version                    int      `json:"version"`
	TotalRowCount              int      `json:"totalRowCount"`
	AccessLevel                string   `json:"accessLevel"`
	EffectiveAttachmentOptions []string `json:"effectiveAttachmentOptions"`
	GanttEnabled               bool     `json:"ganttEnabled"`
	DependenciesEnabled        bool     `json:"dependenciesEnabled"`
	ResourceManagementEnabled  bool     `json:"resourceManagementEnabled"`
	CellImageUploadEnabled     bool     `json:"cellImageUploadEnabled"`
	UserSettings               struct {
		CriticalPathEnabled bool `json:"criticalPathEnabled"`
		DisplaySummaryTasks bool `json:"displaySummaryTasks"`
	} `json:"userSettings"`
	Permalink  string    `json:"permalink"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Columns    []Column  `json:"columns"`
	Rows       []Row     `json:"rows"`
}

//IDToA will return a string representation of the sheetId for easier usage within the SSClient
func (s *Sheet) IDToA() string {
	return strconv.FormatInt(s.ID, 10)
}

//FindValue will search the rows and cols of the sheet looking for a match based on DisplayValue.
//When a value is found, it will return true along with its row and col positions
func (s *Sheet) FindValue(val string) (r *Row, c *Cell, exists bool) {
	for _, r := range s.Rows {
		for _, c := range r.Cells {
			if strings.Compare(c.DisplayValue, val) == 0 {
				//found, remove that from our search list?
				return &r, &c, true
			}
		}
	}
	return nil, nil, false
}

//FindValues will search the rows and cols of the sheet looking for matches based on DisplayValue.
//It will return the items that were not found on the sheet.
func (s *Sheet) FindValues(vals []string) (valsNotFound []string) {
	//convert array to map for easy lookups
	m := make(map[string]bool)
	for _, v := range vals {
		m[strings.ToLower(v)] = false
	}

	for _, r := range s.Rows {
		for _, c := range r.Cells {
			searchVal := strings.ToLower(c.DisplayValue)
			if _, exists := m[searchVal]; exists {
				delete(m, searchVal)
			}
		}
		if len(m) == 0 {
			break //exit if everything was found
		}
	}

	//at the end the map only contains not found items, build array to send back
	retList := make([]string, len(m))
	i := 0
	for k := range m {
		retList[i] = k
		i++
	}

	return retList
}

//Column is a SmartSheet column
type Column struct {
	ID      int64    `json:"id"`
	Index   int      `json:"index"`
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Primary bool     `json:"primary,omitempty"`
	Width   int      `json:"width"`
	Options []string `json:"options,omitempty"`
}

//Row is a SmartSheet row
type Row struct {
	ID             int64      `json:"id,omitempty"`
	RowNumber      int        `json:"rowNumber,omitempty"`
	Expanded       bool       `json:"expanded,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	ModifiedAt     *time.Time `json:"modifiedAt,omitempty"`
	Cells          []Cell     `json:"cells"`
	ParentID       int64      `json:"parentId,omitempty"`
	InCriticalPath bool       `json:"inCriticalPath,omitempty"`
	Locked         bool       `json:"locked,omitempty"`

	//row attributes for location, etc
	ToTop    bool  `json:"toTop,omitempty"`
	ToBottom bool  `json:"toBottom,omitempty"`
	Ident    int16 `json:"ident,omitempty"`
	Outdent  int16 `json:"outdent,omitempty"`

	//Above will never be populated on resposnes, but can be used on Requests
	Above     bool  `json:"above,omitempty"`
	SiblingID int64 `json:"siblingId,omitempty"`
}
