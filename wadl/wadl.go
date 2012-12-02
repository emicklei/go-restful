// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license 
// that can be found in the LICENSE file.

// Package wadl implements the structure for representing a REST-style Webservice API in WADL
// https://wikis.oracle.com/display/Jersey/WADL
package wadl

import (
	"encoding/xml"
)

type Application struct {
	XMLName   xml.Name `xml:"application"`
	Doc       string   `xml:"doc"`
	Resources Resources
}
type Resources struct {
	XMLName   xml.Name `xml:"resources"`
	Base      string   `xml:"base,attr"`
	Resources []Resource
}

func (self *Resources) AddResource(resource Resource) {
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
	Id       string   `xml:"id,attr"`
	Response Response
	Request  Request
	Doc      string `xml:"doc"`
}

type Request struct {
	XMLName        xml.Name `xml:"request"`
	Representation []Representation
}

type Response struct {
	XMLName        xml.Name `xml:"response"`
	Representation []Representation
}

func (self *Request) AddRepresentation(repres Representation) {
	self.Representation = append(self.Representation, repres)
}

func (self *Response) AddRepresentation(repres Representation) {
	self.Representation = append(self.Representation, repres)
}

type Representation struct {
	XMLName   xml.Name `xml:"representation"`
	MediaType string   `xml:"mediaType,attr"`
}
