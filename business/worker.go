package business

import (
	"reporter/util"
)

func collectXmlReportUrl(path string, baseUrl string) ([]string, error) {
	collectedXml, err := util.GetFilesByPostfix(path, ".xml")
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0)
	for _, xml := range collectedXml {
		urls = append(urls, baseUrl+xml)
	}
	return urls, nil
}

func crawlReportData(urls []string) ([]string, error) {

	return nil, nil
}
