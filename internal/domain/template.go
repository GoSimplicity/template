package domain

import (
	"fmt"
	"github.com/GoSimplicity/template/internal/errs"
)

type Template struct {
	ID    int64 `json:"id"`
	BizID int64 `json:"bizId"` // 业务唯一标识
	// ...
}

func (t *Template) Validate() error {
	if t.BizID <= 0 {
		return fmt.Errorf("%w: BizID = %d", errs.ErrInvalidParameter, t.BizID)
	}

	if t.ID <= 0 {
		return fmt.Errorf("%w: ID = %d", errs.ErrInvalidParameter, t.ID)
	}

	return nil
}
