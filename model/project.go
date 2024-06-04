package model

import "encoding/xml"

type Project struct {
	XMLName  xml.Name `xml:"project"`
	Name     string   `xml:"name,attr"`
	Property []struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	} `xml:"property"`
}

// factory from xml
func NewProjectFromXml(data []byte) (*Project, error) {
	var project Project
	err := xml.Unmarshal(data, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}
