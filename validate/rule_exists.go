package validate

import (
	"gorm.io/gorm"
)

// check if resource exists (where business_id = ?)
type ruleExists struct {
	table string
	id    any
	err   error
	do    *bool
	*HasFilter
}

// specifies When to validate
// if When is not specified, will validate by default
func (rule ruleExists) When(when bool) ruleExists {
	rule.do = &when
	return rule
}

func (vr ruleExists) Init() bool {
	// skip validation if user specifies when
	return vr.do == nil || *vr.do
}

func (vr ruleExists) CountResults(dbCtx *gorm.DB) error {
	var count int64
	dbCtx = dbCtx.Table(vr.table).Where("id = ?", vr.id)
	vr.ApplyFilter(dbCtx)
	if err := dbCtx.Count(&count).Error; err != nil {
		return err
	}
	if count <= 0 {
		return vr.err
	}

	return nil
}

func NewExistsRule(table string, id any, err error, filter *HasFilter) ruleExists {
	return ruleExists{
		table:     table,
		id:        id,
		HasFilter: filter,
		err:       err,
	}
}
