package class

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/scottdware/go-panos"
)

const(
	xpathNetflow = "/config/shared/server-profile/netflow"
)

type NetflowServer struct {
	XMLName xml.Name `xml:"entry"`

	Name string `xml:"name,attr"`
	Host string `xml:"host"`
	Port int `xml:"port,omitempty"`
}

type TemplateRefresh struct {
	XMLName xml.Name `xml:"template-refresh-rate"`

	Minutes int `xml:"minutes,omitempty"`
	Packets int `xml:"packets,omitempty"`
}

type ListOfServer struct {
	XMLName xml.Name `xml:"server"`
	NetflowServer []NetflowServer `xml:"entry"`
}

type NetflowProfile struct {
	XMLName xml.Name `xml:"entry"`

	Xpath string `xml:"-"`
	Id string `xml:"-"`
	Name string `xml:"name,attr"`
	Server ListOfServer `xml:"server"`
	ActiveTimeout int `xml:"active-timeout,omitempty"`
	TemplateRefresh TemplateRefresh `xml:"template-refresh-rate,omitempty"`
}

type ResultListOfNetflowProfiles struct {
	Profiles []NetflowProfile `xml:"netflow>entry"`
}

type XMLNetflowProfile struct {
	Profile NetflowProfile `xml:"entry"`
}

func (n *NetflowProfile)Flatten() []interface{} {
	c := make(map[string]interface{}, 0)
	c["name"] = n.Name
	c["active-timeout"] = n.ActiveTimeout
	c["template-refresh-rate"] = n.TemplateRefresh
	c["server"] = n.Server
	return []interface{}{c}
}

func (n *NetflowProfile) ToXML() string {
	result, err := xml.Marshal(n)
	if(err != nil) {
		fmt.Println("Error",err)
		return ""
	}
	return string(result)
}

func (n *NetflowProfile) ToString() string {
	return n.Xpath + "\n" + n.ToXML()
}

func (n *NetflowProfile) GetXpath() string {
	return n.Xpath
}

func (n *NetflowProfile) GetName() string {
	return n.Name
}

func (n *NetflowProfile) Add(pan *panos.PaloAlto) error {
	err := pan.XpathConfig("edit",fmt.Sprintf(xpathNetflow + "/entry[@name='%s']",n.Name),n.ToXML())
	return err
}

func (n *NetflowProfile) Delete(pan *panos.PaloAlto) error {
	err := pan.XpathConfig("delete",fmt.Sprintf(xpathNetflow + "/entry[@name='%s']",n.Name),"")
	return err
}

func (n *NetflowProfile) Edit(pan *panos.PaloAlto) error {
	if(n.Id != n.Name) {
		fmt.Println("Changing name of object")
		err := n.Delete(pan)
		if(err != nil) {
			return err
		}
		n.Id = n.Name
	}
	return n.Add(pan)
}

/*
type NetflowServer struct {
	XMLName xml.Name `xml:"entry"`

	Name string `xml:"name,attr"`
	Host string `xml:"host"`
	Port int `xml:"port,omitempty"`
}

type TemplateRefresh struct {
	XMLName xml.Name `xml:"entry"`

	Minutes int `xml:"minutes,omitempty"`
	Packets int `xml:"packets,omitempty"`
}

type NetflowProfile struct {
	XMLName xml.Name `xml:"entry"`

	Xpath string `xml:"-"`
	Id string `xml:"-"`
	Name string `xml:"name,attr"`
	Server []NetflowServer `xml:"server"`
	ActiveTimeout int `xml:"active-timeout,omitempty"`
	TemplateRefresh TemplateRefresh `xml:"template-refresh-rate,omitempty"`
}
*/

func CreateNetflowClass(name string, server []NetflowServer, TemplateRefresh TemplateRefresh, activeTimeout int) NetflowProfile {
	return NetflowProfile{ 
		Name: name , 
		Xpath: xpathNetflow, 
		Id: name,
		Server: ListOfServer{NetflowServer: server},
		TemplateRefresh: TemplateRefresh,
		ActiveTimeout: activeTimeout,
	}
}

func SearchNetflow(name string, pan *panos.PaloAlto) (NetflowProfile,error) {
	response, err := pan.XpathGetConfig("candidate", xpathNetflow + "/entry[@name='" + name + "']")
	if(err != nil) {
		return NetflowProfile{},err
	}
	parsedResponse := Response[Result[NetflowProfile]]{}
	err = xml.Unmarshal([]byte(response), &parsedResponse)
	if parsedResponse.Response.Entry.Name == "" {
		return NetflowProfile{}, errors.New("Netflow not found")
	}
	return parsedResponse.Response.Entry, err
}

func CheckIfNetflowExist(name string, pan *panos.PaloAlto) bool {
	_, err := SearchNetflow(name,pan)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return err == nil
}