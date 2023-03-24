package class

import (
	"fmt"

	"github.com/scottdware/go-panos"
)

type FWObj interface {
	ToXML() string
	GetXpath() string
	GetName() string
}

func FWAddObject(pan *panos.PaloAlto, obj FWObj) error {
	err := pan.XpathConfig("edit",fmt.Sprintf(obj.GetXpath() + "/entry[@name='%s']"),obj.ToXML())
	if(err != nil) {
		return err
	}
	return nil
}

func Delete(pan *panos.PaloAlto, obj FWObj) error {
	err := pan.XpathConfig("delete",fmt.Sprintf(obj.GetXpath() + "/entry[@name='%s']",obj.GetName()), obj.ToXML())
	if(err != nil) {
		return err
	}
	return nil
}