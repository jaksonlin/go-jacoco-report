package util_test

import (
	"os"
	"reporter/util"
	"testing"
)

func TestDownloadTestReport(t *testing.T) {
	url := "https://www.baidu.com"
	body, err := util.DownloadWebPage(url)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fp, err := os.Create("test.html")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer fp.Close()
	fp.Write(body)

}
