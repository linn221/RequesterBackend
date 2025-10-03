package validate

import (
	"gorm.io/gorm"
)

type noResultRule struct {
	statusCode *int
	table      string
	err        error
	do         *bool
	*HasFilter
}

// specifies When to validate
// if When is not specified, will validate by default
func (rule noResultRule) When(when bool) noResultRule {
	rule.do = &when
	return rule
}

func (rule noResultRule) OverrideStatusCode(i int) noResultRule {
	rule.statusCode = &i
	return rule
}

func (vr noResultRule) Init() bool {
	// skip validation if user specifies when
	return vr.do == nil || *vr.do
}

func (vr noResultRule) CountResults(dbCtx *gorm.DB) error {
	var count int64
	dbCtx = dbCtx.Table(vr.table)
	vr.ApplyFilter(dbCtx)
	if err := dbCtx.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return vr.err
	}

	return nil
}

func NewNoResultRule(table string, err error, filter *HasFilter) noResultRule {
	return noResultRule{
		table:     table,
		err:       err,
		HasFilter: filter,
	}
}
