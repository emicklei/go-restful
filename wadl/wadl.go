package wadl
/*
https://wikis.oracle.com/display/Jersey/WADL
*/
import (
	"encoding/xml"
)

type Application struct {
	XMLName   xml.Name `xml:"application"`
	Doc       string   `xml:"doc"`
	Resources []Resource
}

func (self *Application) AddResource(resource Resource) {
	self.Resources = append(self.Resources, resource)
}

type Resource struct {
	XMLName xml.Name `xml:"resource"`
	Method  Method
	Path    string `xml:"path,attr"`
}

type Method struct {
	XMLName  xml.Name `xml:"method"`
	Name     string   `xml:"name,attr"`
//	Id       string   `xml:"id,attr"`
	Response Response
}

type Response struct {
	XMLName        xml.Name `xml:"response"`
	Representation Representation
}

type Representation struct {
	XMLName   xml.Name `xml:"representation"`
	MediaType string   `xml:"mediaType,attr"`
}
