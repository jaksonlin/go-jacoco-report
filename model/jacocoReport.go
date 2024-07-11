package model

import (
	"bytes"
	"container/list"
	"fmt"
	"regexp"
	"reporter/util"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type JacocoReport struct {
	ReportTitle   string                `json:"report_title"`
	ReportDate    string                `json:"report_date"`
	ReportContent []*JacocoReportRecord `json:"report_content"`
}

type JacocoReportRecord struct {
	Element             string        `json:"element"`
	Href                string        `json:"href"`
	InstructionsCov     string        `json:"instruct_coverage"`
	BranchCov           string        `json:"branch_coverage"`
	MissedBranch        string        `json:"missed_branch"`
	Complexity          string        `json:"complexity"`
	MissedMethod        string        `json:"missed_method"`
	NumberOfMethod      string        `json:"number_of_method"`
	MissedLine          string        `json:"missed_line"`
	NumberOfLine        string        `json:"number_of_line"`
	MissedClass         string        `json:"missed_class"`
	NumberOfClass       string        `json:"number_of_class"`
	SubReport           *JacocoReport `json:"sub_report"`
	MissedBranchCount   uint64
	ComplexityCount     uint64
	MissedMethodCount   uint64
	NumberOfMethodCount uint64
	MissedLineCount     uint64
	NumberOfLineCount   uint64
	MissedClassCount    uint64
	NumberOfClassCount  uint64
}

func (j *JacocoReportRecord) ConvertNumeric() error {
	c, err := strconv.ParseUint(j.MissedBranch, 10, 64)
	if err != nil {
		return err
	}
	j.MissedBranchCount = c
	c, err = strconv.ParseUint(j.Complexity, 10, 64)
	if err != nil {
		return err
	}
	j.ComplexityCount = c
	c, err = strconv.ParseUint(j.MissedMethod, 10, 64)
	if err != nil {
		return err
	}
	j.MissedMethodCount = c
	c, err = strconv.ParseUint(j.NumberOfMethod, 10, 64)
	if err != nil {
		return err
	}
	j.NumberOfMethodCount = c
	c, err = strconv.ParseUint(j.MissedLine, 10, 64)
	if err != nil {
		return err
	}
	j.MissedLineCount = c
	c, err = strconv.ParseUint(j.NumberOfLine, 10, 64)
	if err != nil {
		return err
	}
	j.NumberOfLineCount = c
	c, err = strconv.ParseUint(j.MissedClass, 10, 64)
	if err != nil {
		return err
	}
	j.MissedClassCount = c
	c, err = strconv.ParseUint(j.NumberOfClass, 10, 64)
	if err != nil {
		return err
	}
	j.NumberOfClassCount = c
	return nil
}

var titleRegex *regexp.Regexp = regexp.MustCompile(`<h1>(.*)</h1>`)

func NewJacocoReportFromURL(url string) (*JacocoReport, error) {
	// download web page
	body, err := util.DownloadWebPage(url)
	if err != nil {
		return nil, err
	}

	return NewJacocoReportFromHtmlBytes(body)
}

func (j *JacocoReport) CrawSubReports(url string) {
	jobList := list.New()
	// work recorder
	type crawlSettings struct {
		urlTowork string
		report    *JacocoReportRecord
	}
	// init jobList with search on current direct child
	for index := range j.ReportContent {
		if len(j.ReportContent[index].Href) > 0 && !strings.Contains(j.ReportContent[index].Href, "#") {
			jobList.PushBack(&crawlSettings{
				urlTowork: url + "/" + j.ReportContent[index].Href,
				report:    j.ReportContent[index],
			})
		}
	}

	errCount := 0
	for jobList.Len() > 0 {
		workItem := jobList.Front().Value.(*crawlSettings)
		jobList.Remove(jobList.Front())

		subReport, err := NewJacocoReportFromURL(workItem.urlTowork)
		if err != nil {
			errCount += 1
			fmt.Printf("download report for %s, failed, link: %s\n", subReport.ReportTitle, workItem.urlTowork)
			continue
		}

		workItem.report.SubReport = subReport
		// expand search on child
		for index := range subReport.ReportContent {
			if len(j.ReportContent[index].Href) > 0 && !strings.Contains(j.ReportContent[index].Href, "#") {
				urlBase := workItem.urlTowork // current work item url as next level base url
				if strings.Contains(urlBase, "index.html") {
					urlBase = strings.TrimRight(urlBase, "index.html")
				}
				jobList.PushBack(&crawlSettings{
					urlTowork: urlBase + j.ReportContent[index].Href,
					report:    j.ReportContent[index],
				})
			}
		}
	}
	fmt.Println("total failure: ", errCount)
}

func NewJacocoReportFromHtmlBytes(body []byte) (*JacocoReport, error) {
	title := titleRegex.FindAllSubmatch(body, -1)
	if len(title) == 0 {
		return nil, fmt.Errorf("title not found")
	}
	if len(title) > 1 {
		return nil, fmt.Errorf("multiple titles found")
	}

	report := &JacocoReport{
		ReportTitle: string(title[0][1]),
		ReportDate:  time.Now().Format("2006-01-02 15:04:05"),
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	table := util.GetHtmlTable(doc)
	if table == nil {
		return nil, fmt.Errorf("table not found")
	}

	report.ReportContent, err = createJacocoReportRecords(table)
	if err != nil {
		return nil, err
	}
	return report, nil

}

func createJacocoReportRecords(table *html.Node) ([]*JacocoReportRecord, error) {
	var records []*JacocoReportRecord = make([]*JacocoReportRecord, 0)
	recordIndex := 0
	tableColIndex := 0
	err := util.TraverseJacocoHtmlTable(table, func(data string, totalColumnCount int) error {
		contentDetails := strings.Split(data, ";")
		text := ""
		href := ""
		for _, contentDetail := range contentDetails {
			contentDetailParts := strings.Split(contentDetail, "=")
			if len(contentDetailParts) != 2 {
				return fmt.Errorf("invalid content detail: %s", contentDetail)
			}
			if contentDetailParts[0] == "text" {
				text = contentDetailParts[1]
			} else if contentDetailParts[0] == "href" {
				href = contentDetailParts[1]
			}
		}
		switch tableColIndex {
		case 0:
			records = append(records, &JacocoReportRecord{
				Element: text,
				Href:    href,
			})
		case 2:
			records[recordIndex].InstructionsCov = text
		case 4:
			records[recordIndex].BranchCov = text
		case 5:
			records[recordIndex].MissedBranch = text
		case 6:
			records[recordIndex].Complexity = text
		case 7:
			records[recordIndex].MissedLine = text
		case 8:
			records[recordIndex].NumberOfLine = text
		case 9:
			records[recordIndex].MissedMethod = text
		case 10:
			records[recordIndex].NumberOfMethod = text
		case 11:
			records[recordIndex].MissedClass = text
		case 12:
			records[recordIndex].NumberOfClass = text
		}
		tableColIndex += 1
		if tableColIndex == totalColumnCount {
			recordIndex += 1
			tableColIndex = 0
			err := records[recordIndex].ConvertNumeric()
			if err != nil {
				fmt.Printf("convert numeric failed: %s\n", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (j *JacocoReport) GetHighComplexityFunctions(thredshold uint64) []*JacocoReportRecord {
	var highComplexityFunctions []*JacocoReportRecord = make([]*JacocoReportRecord, 0)
	searchList := list.New()
	for _, record := range j.ReportContent {
		searchList.PushBack(record)
	}
	for searchList.Len() > 0 {
		searchItem := searchList.Front().Value.(*JacocoReportRecord)
		searchList.Remove(searchList.Front())
		if strings.Contains(searchItem.Href, ".java.html#") && searchItem.ComplexityCount > thredshold {
			highComplexityFunctions = append(highComplexityFunctions, searchItem)
		}
		if searchItem.SubReport != nil {
			for _, record := range searchItem.SubReport.ReportContent {
				searchList.PushBack(record)
			}
		}
	}
	return highComplexityFunctions
}

func (j *JacocoReport) GetLowCoverageClassess(thredshold float64) []*JacocoReportRecord {
	var lowCoverageClassess []*JacocoReportRecord = make([]*JacocoReportRecord, 0)
	searchList := list.New()
	for _, record := range j.ReportContent {
		searchList.PushBack(record)
	}
	for searchList.Len() > 0 {
		searchItem := searchList.Front().Value.(*JacocoReportRecord)
		searchList.Remove(searchList.Front())
		if searchItem.SubReport != nil && len(searchItem.SubReport.ReportContent) > 0 {

			if strings.Contains(searchItem.SubReport.ReportContent[0].Href, ".java.html#") {
				if searchItem.NumberOfLineCount > 0 {
					coverage := 1 - float64(searchItem.MissedLineCount)/float64(searchItem.NumberOfLineCount)
					if coverage < thredshold {
						lowCoverageClassess = append(lowCoverageClassess, searchItem)
					}
				}
			} else {
				for _, record := range searchItem.SubReport.ReportContent {
					searchList.PushBack(record)
				}
			}

		}

	}
	return lowCoverageClassess
}
