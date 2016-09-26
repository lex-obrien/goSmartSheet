package goSmartSheet

import (
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
	for k := range m {
		retList = append(retList, k)
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
	ID         int64     `json:"id"`
	RowNumber  int       `json:"rowNumber"`
	Expanded   bool      `json:"expanded"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Cells      []Cell    `json:"cells"`
	ParentID   int64     `json:"parentId,omitempty"`
	SiblingID  int64     `json:"siblingId,omitempty"`
	//row attributes for location, etc
	//TODO: put these back in...
	// ToTop bool
	// ToBottom bool
	// Above bool
	// Ident int16
	// Outdent int16
}
