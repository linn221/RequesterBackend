package validate

import (
	"gorm.io/gorm"
)

// check if slice of resource id exists (where business_id IN ?)
type RuleMassExists[ID comparable] struct {
	Table         string
	Ids           []ID
	Err           error
	NoDuplicateID bool
	*HasFilter
}

func (r RuleMassExists[ID]) Init() bool {
	return len(r.Ids) > 0
}

func (r RuleMassExists[ID]) CountResults(dbCtx *gorm.DB) error {
	var count int64
	uniqIds := UniqueSlice(r.Ids)
	dbCtx = dbCtx.Table(r.Table).Where("id IN ?", uniqIds)
	err := dbCtx.Count(&count).Error
	if err != nil {
		return err
	}
	if count != int64(len(uniqIds)) {
		return r.Err
	}

	return nil
}
