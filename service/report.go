package service

import (
	"context"
	"fmt"
)

type RdsRefreshReport struct {
	TotalCount     int64 `json:"total_count"`
	ProcessedCount int64 `json:"processed_count"`

	SatisfiedErrCount int64 `json:"satisfied_err_count"`
	RefreshErrCount   int64 `json:"refresh_err_count"`

	UnsatisfiedCount     int64 `json:"unsatisfied_count"`
	UnindeedRefreshCount int64 `json:"unindeed_refresh_count"`

	IndeedRefreshCount int64 `json:"indeed_refresh_count"`

	ErrDetailUrl string `json:"err_detail_url"`

	LastProcessedID int64 `json:"last_processed_id"`
}

type ReportType string

const (
	FinalReport   ReportType = "Final Report"
	SegmentReport ReportType = "Segment Report"
)

func AppendReportRow(rowData []string) {

}

func DoNotification(ctx context.Context, operator string,
	title string, reportType ReportType, report *RdsRefreshReport, err error) error {

	title = title + string(reportType)

	if err != nil {
		fmt.Printf("【data_manipulation_service】DoNotification occur error: %s\n", err.Error())
		return nil

	}

	if report == nil {
		return nil
	}

	var content string
	content += fmt.Sprintf("| ---------------------- |    \n")
	content += fmt.Sprintf("| Operator: %s \n", operator)
	content += fmt.Sprintf("| ---------------------- |    \n")
	content += fmt.Sprintf("| Total Count            | %d \n", report.ProcessedCount)
	content += fmt.Sprintf("| Indeed Refresh Count   | %d \n", report.IndeedRefreshCount)
	content += fmt.Sprintf("| Unsatisfied Count      | %d \n", report.UnsatisfiedCount)
	content += fmt.Sprintf("| Unindeed Refresh Count | %d \n", report.UnindeedRefreshCount)
	content += fmt.Sprintf("| ---------------------- |    \n")
	content += fmt.Sprintf("| Satisfied Err Count    | %d \n", report.SatisfiedErrCount)
	content += fmt.Sprintf("| Refresh Err Count      | %d \n", report.RefreshErrCount)
	content += fmt.Sprintf("| ---------------------- |    \n")
	content += fmt.Sprintf("| Error Details: %s \n", report.ErrDetailUrl)

	return Notify(title, content, operator)

}
