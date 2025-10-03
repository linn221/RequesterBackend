package validate

import (
	"errors"

	"gorm.io/gorm"
)

type ruleUnique struct {
	table    string
	err      error
	column   string
	value    any
	exceptId int
	do       *bool
	message  string

	*HasFilter
}

func (rule ruleUnique) When(cond bool) ruleUnique {
	rule.do = &cond
	return rule
}

// func (rule ruleUnique) Filter(cond string, values ...any) ruleUnique {
// 	rule.HasFilter = &HasFilter{
// 		Cond:         cond,
// 		FilterValues: values,
// 	}
// 	return rule
// }

func (rule ruleUnique) Init() bool {

	if rule.do != nil && !*rule.do {
		return false
	}
	return true
}

func (r ruleUnique) CountResults(dbCtx *gorm.DB) error {
	var count int64

	query := "`" + r.column + "`" + " = ?"
	dbCtx = dbCtx.Table(r.table).Where(query, r.value)
	if r.exceptId > 0 {
		dbCtx.Where("id != ?", r.exceptId)
	}
	r.ApplyFilter(dbCtx)
	err := dbCtx.Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(r.message)
	}

	return nil
}

func (r ruleUnique) Except(id int) ruleUnique {
	r.exceptId = id
	return r
}

func (r ruleUnique) Say(message string) ruleUnique {
	r.message = message
	return r
}

func NewUniqueRule(table string, column string, value any, filter *HasFilter) ruleUnique {
	// var v T
	return ruleUnique{
		table:     table,
		column:    column,
		value:     value,
		message:   "duplicate record somewhere",
		HasFilter: filter,
	}
}
