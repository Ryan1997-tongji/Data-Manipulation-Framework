// package impl
// @author: chenzhewei.97
// @create date: 2025/6/10
package impl

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type TestRefresher struct {
}

type TestRefresherInput struct {
	FieldX string
	FieldY int64
}

func (r *TestRefresher) GetFilterCondition(ctx context.Context, input interface{}) (condition *gorm.DB, err error) {
	// your db name and filter condition
	return condition, nil
}

func (r *TestRefresher) IsOneSatisfied(ctx context.Context, input interface{}, id int64) (isSatisfied bool, err error) {
	// your biz logic
	return isSatisfied, nil
}

func (r *TestRefresher) RefreshOne(ctx context.Context, input interface{}, id int64) (isIndeedRefreshed bool, err error) {
	fmt.Println("TestRefresher")
	// your biz logic
	return true, nil
}

func (r *TestRefresher) GetInputInfo(ctx context.Context, input interface{}, id int64) *TestRefresherInput {
	inputInfo, ok := input.(*TestRefresherInput)
	if ok && inputInfo != nil {
		return inputInfo
	}

	return nil
}
