// package service
// @author: chenzhewei.97
// @create date: 2025/5/28
package service

import (
	"context"
	"gorm.io/gorm"
)

type Refresher interface {

	// GetFilterCondition
	GetFilterCondition(ctx context.Context, input interface{},
	) (condition *gorm.DB, err error)

	// IsOneSatisfied
	IsOneSatisfied(ctx context.Context, input interface{},
		id int64) (isSatisfied bool, err error)

	// RefreshOne
	RefreshOne(ctx context.Context, input interface{},
		id int64) (isIndeedRefreshed bool, err error)
}
