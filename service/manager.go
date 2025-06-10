package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ryan1997-tongji/Data-Manipulation-Framework/integration"
	"github.com/Ryan1997-tongji/Data-Manipulation-Framework/utils"
	"gorm.io/gorm"
	"time"
)

var reportSegment int64 = 100000

func DoRefresh(ctx context.Context,
	taskTitle string, operator string,
	input interface{}, refresher Refresher) {
	fmt.Printf("【data_manipulation_service】DoRefresh【input】%s\n", utils.SimpleJson(input))
	utils.SafeGo(func() {
		doRefreshWithNotification(ctx,
			taskTitle, operator,
			input, refresher)
	})

}

func doRefreshWithNotification(ctx context.Context,
	taskTitle string, operator string,
	input interface{}, refresher Refresher) {

	var err error

	var report *RdsRefreshReport
	report, err = Exec(ctx, taskTitle, operator, input, refresher)
	fmt.Printf("【data_manipulation_service】DoRefresh【finished】【report】%s\n", utils.SimpleJson(report))

	err = DoNotification(ctx, operator, taskTitle, FinalReport, report, err)
	if err != nil {
		fmt.Printf("【data_manipulation_service】DoReport occur error: %s\n", err.Error())
	}

}

func Exec(ctx context.Context, taskTitle string,
	operator string, input interface{}, refresher Refresher) (report *RdsRefreshReport, err error) {

	// init
	report = &RdsRefreshReport{}

	// 2. create report
	report.ErrDetailUrl = integration.CreateReport(operator)

	var batchSize int64 = 200

	var curBeginId int64 = 0
	var nextBeginId int64 = 0

	var curQueryErrCount int64 = 0

	for {
		// 1. batch query
		var ids []int64
		ids, nextBeginId, err = BatchQueryIds(ctx, input, refresher, batchSize, curBeginId)
		if err != nil {

			curQueryErrCount++
			fmt.Printf(
				"【data_manipulation_service】DoRefresh【err occur】BatchQueryIds【times】%d【beginId】%d【error】%s",
				curQueryErrCount, curBeginId, err.Error())

			if curQueryErrCount >= 3 {
				// more than 3 times error
				return report, errors.New(fmt.Sprintf("【data_manipulation_service】DoRefresh【fail】BatchQueryIds Exceed 3 times【beginId】%d【error】%s\n", curBeginId, err.Error()))
			} else {
				time.Sleep(100 * time.Millisecond)
				continue
			}

		} else {

			// clear the error count
			curQueryErrCount = 0
			// judge whether to break
			if len(ids) <= 0 || curBeginId == nextBeginId {
				return report, nil
			}
			// set curBeginId=nextBeginId
			curBeginId = nextBeginId

		}

		// 2. statistics
		report.ProcessedCount += int64(len(ids))
		report.LastProcessedID = ids[0]
		fmt.Printf(
			"【data_manipulation_service】DoRefresh【processing】now at %d",
			report.ProcessedCount)
		// 3. if meet segment report condition, do segment report
		//if report.ProcessedCount > 0 && report.ProcessedCount%reportSegment == 0 {
		//	err = DoNotification(ctx, operator, taskTitle, SegmentReport, report, nil)
		//	if err != nil {
		//		fmt.Printf("【data_manipulation_service】DoReport occur error: %s\n", err.Error())
		//	}
		//}

		// 4. refresh each one
		for _, _id := range ids {

			var id int64
			id = _id

			// 4.1 判断是否符合清洗条件
			var isSatisfied bool
			isSatisfied, err = refresher.IsOneSatisfied(ctx, input, id)
			if err != nil {
				report.SatisfiedErrCount++
				fmt.Printf(
					"【data_manipulation_service】DoRefresh【err occur】IsOneSatisfied【id】%d【input】%s【error】%s",
					id, utils.SimpleJson(input), err.Error())
				continue
			} else {
				if !isSatisfied {
					report.UnsatisfiedCount++
					continue
				}
			}

			// 4.2 do refresh
			var isIndeedRefreshed bool
			isIndeedRefreshed, err = refresher.RefreshOne(ctx, input, id)
			if err != nil {
				report.RefreshErrCount++
				fmt.Printf(
					"【data_manipulation_service】DoRefresh【err occur】RefreshOne【id】%d【input】%s【error】%s",
					id, utils.SimpleJson(input), err.Error())
				AppendReportRow([]string{err.Error()})
			} else {
				if !isIndeedRefreshed {
					report.UnindeedRefreshCount++
				} else {
					report.IndeedRefreshCount++
				}
			}

		}

	}

}

func BatchQueryIds(ctx context.Context, input interface{}, refresher Refresher,
	batchSize int64, beginId int64) (ids []int64, nextBeginId int64, err error) {

	var dbr *gorm.DB
	dbr, err = refresher.GetFilterCondition(ctx, input)
	if err != nil {
		return nil, beginId, err
	}

	err = dbr.
		Where("id > ?", beginId).
		Select("id").
		Order("id ASC").
		Limit(int(batchSize)).
		Find(&ids).Error
	if err != nil {
		return nil, beginId, err
	}

	if len(ids) <= 0 {
		return nil, beginId, nil
	}

	return ids, ids[len(ids)-1], nil

}
