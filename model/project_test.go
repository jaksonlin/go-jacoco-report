package model_test

import (
	"reporter/model"
	"testing"
)

func TestWorker(t *testing.T) {
	// test worker

	xmlPayload := `<?xml version="1.0" encoding="UTF-8"?>
	<project name="amc-amccsa-biz-trunk-jacoco" default="all-report" basedir="." xmlns:jacoco="antlib:org.jacoco.ant">
    <taskdef uri="antlib:org.jacoco.ant" resource="org/jacoco/ant/antlib.xml">  
        <classpath path="${basedir}/lib/jacocoant.jar"/>
    </taskdef>
    <property name="report_number" value="master-amc-amccsa-unittest-report" />
	</project>
	</xml>`

	project, err := model.NewProjectFromXml([]byte(xmlPayload))
	if err != nil {
		t.Errorf("Error parsing xml: %v", err)
	}
	if project.Name != "amc-amccsa-biz-trunk-jacoco" {
		t.Errorf("Expected project name to be 'amc-amccsa-biz-trunk-jacoco', got '%s'", project.Name)
	}

	// test worker

}
