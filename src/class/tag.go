package class

import (
	"encoding/xml"
	"errors"
	//"errors"
	//"errors"
	"fmt"

	"github.com/scottdware/go-panos"
)

const(
	xpathTag = "/config/devices/entry/vsys/entry/tag"
)

type Tag struct {
	XMLName xml.Name `xml:"entry"`

	Xpath string `xml:"-"`
	Id string `xml:"-"`
	Name     string `xml:"name,attr"`
	Color    string `xml:"color,omitempty"`
	Comments string `xml:"comments,omitempty"`
}


type ResponseTag struct {
	XMLName  xml.Name `xml:"entry"`
	Xpath    string   `xml:"-"`
	Id       string   `xml:"-"`
	Name     string   `xml:"name,attr"`
	Color    string   `xml:"color,omitempty"`
	Comments string   `xml:"comments,omitempty"`
}

type ResultListOfTags struct {
	Tags []Tag `xml:"tag>entry"`
}

type XMLTag struct {
	Tag Tag `xml:"entry"`
}

func (t *Tag)Flatten() []interface{} {
	c := make(map[string]interface{}, 0)
	c["name"] = t.Name
	c["color"] = t.Color
	c["comments"] = t.Comments
	return []interface{}{c}
}

func (t *Tag) ToXML() string {
	result, err := xml.Marshal(t)
	if(err != nil) {
		fmt.Println("Error",err)
		return ""
	}
	return string(result)
}

func (t *Tag) ToString() string {
	return t.Xpath + "\n" + t.ToXML()
}

func (t *Tag) GetXpath() string {
	return t.Xpath
}

func (t *Tag) GetName() string {
	return t.Name
}

func (t *Tag) Add(pan *panos.PaloAlto) error {
	err := pan.XpathConfig("edit",fmt.Sprintf(xpathTag + "/entry[@name='%s']",t.Name),t.ToXML())
	if(err != nil) {
		fmt.Println(t.Xpath + "/entry[@name='%s']",t.Name)
		fmt.Println(t.ToXML())
		return err
	}
	return nil
}

func (t *Tag) Edit(pan *panos.PaloAlto) error {
	if(t.Id != t.Name) {
		err := t.Delete(pan)
		if(err != nil) {
			return err
		}
		t.Id = t.Name
		t.Add(pan)
	}
	err := pan.XpathConfig("edit",fmt.Sprintf(xpathTag + "/entry[@name='%s']",t.Name),t.ToXML())
	if(err != nil) {
		fmt.Printf(xpathTag + "/entry[@name='%s']\n",t.Name)
		fmt.Println(t.ToXML())
		return err
	}
	return nil
}

func (t *Tag) Delete(pan *panos.PaloAlto) error {
	err := pan.XpathConfig("delete",fmt.Sprintf(xpathTag + "/entry[@name='%s']",t.Id),t.ToXML())
	if(err != nil) {
		return err
	}
	return nil
}

func SearchTag(name string, pan *panos.PaloAlto) (Tag, error) {
	result, err := pan.XpathGetConfig("candidate",xpathTag + "/entry[@name='" + name + "']")
	if(err != nil) {
		fmt.Println("Error",err)
		return Tag{}, err
	}
	ResponseTag := Response[Result[Tag]]{}
	err = xml.Unmarshal([]byte(result), &ResponseTag)
	if(err != nil) {
		fmt.Println("Error",err)
		return Tag{}, err
	}
	if(ResponseTag.Response.Entry.Name == "") {
		return Tag{}, errors.New("Tag not found")
	}
	tag := CreateTagClass(ResponseTag.Response.Entry.Name,ResponseTag.Response.Entry.Color,ResponseTag.Response.Entry.Comments)
	return tag, nil
}

func CheckIfTagExist(name string, pan *panos.PaloAlto) bool {
	_, err := SearchTag(name, pan)
	return err == nil
}

func CreateTagClass(name string, color string, comments string) Tag {
	return Tag{
		Xpath: xpathTag,
		Name: name,
		Color: color,
		Comments: comments,
		Id: name,
	}
}