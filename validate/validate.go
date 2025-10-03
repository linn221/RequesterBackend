package validate

import (
	"gorm.io/gorm"
)

type Rule interface {
	Init() bool
	CountResults(*gorm.DB) error
}

func Validate(db *gorm.DB, rules ...Rule) error {
	for _, rule := range rules {
		if ok := rule.Init(); !ok {
			continue
		}
		err := rule.CountResults(db)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateInBatch(db *gorm.DB, rules ...Rule) []error {
	errors := make([]error, 0)
	for _, rule := range rules {
		if ok := rule.Init(); !ok {
			continue
		}
		err := rule.CountResults(db)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

type HasFilter struct {
	Cond         string
	FilterValues []any
}

func (f *HasFilter) ApplyFilter(dbCtx *gorm.DB) {
	if f != nil {
		dbCtx.Where(f.Cond, f.FilterValues...)
	}
}

func NewFilter(cond string, values ...any) *HasFilter {
	return &HasFilter{
		Cond:         cond,
		FilterValues: values,
	}
}

func NewShopFilter(shopId int) *HasFilter {
	return &HasFilter{
		Cond:         "shop_id = ?",
		FilterValues: []any{shopId},
	}
}
