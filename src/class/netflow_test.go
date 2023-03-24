package class

import (
	"fmt"
	"os"
	//"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/scottdware/go-panos"
	"github.com/stretchr/testify/assert"
)

func TestNetflow(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Errorf("Error loading .env file")
	}
	assert := assert.New(t)
	// Prepare creds for testing
	hostname := os.Getenv("PANOS_HOSTNAME")
	username := os.Getenv("PANOS_USERNAME")
	password := os.Getenv("PANOS_PASSWORD")

	assert.NotNil(hostname, "Hostname is nil")
	assert.NotNil(username, "Username is nil")
	assert.NotNil(password, "Password is nil")

	assert.NotEqual(hostname, "", "Hostname is empty")
	assert.NotEqual(username, "", "Username is empty")
	assert.NotEqual(password, "", "Password is empty")
	
	creds := &panos.AuthMethod{
		Credentials: []string{username, password},
	}

	// Prepare panos session
	pan, err := panos.NewSession(hostname, creds)
	assert.Nil(err)
	
	// Netflow tests
	servers := []NetflowServer{
		{
			Host: "192.10.10.10",
			Port: 9999,
			Name: "toto",
		},
	}
	templates := TemplateRefresh{
			Minutes: 15,
			Packets: 15,
	}
	testNetflow := CreateNetflowClass("tata",servers,templates, 15)
	assert.Equal(testNetflow.Name,"tata")
	assert.Equal(testNetflow.Server.NetflowServer[0].Host, "192.10.10.10")
	assert.Equal(testNetflow.Server.NetflowServer[0].Port, 9999)
	assert.Equal(testNetflow.Server.NetflowServer[0].Name, "toto")
	assert.Equal(testNetflow.TemplateRefresh.Minutes, 15)
	assert.Equal(testNetflow.TemplateRefresh.Packets, 15)
	assert.Equal(testNetflow.ActiveTimeout, 15)

	// Test Add netflow
	errAdd := testNetflow.Add(pan)
	assert.Nil(errAdd)

	// Test Search netflow
	netflow, searchError := SearchNetflow("tata", pan)
	assert.Nil(searchError)
	assert.Equal(netflow.Name, "tata")
	assert.Equal(netflow.Server.NetflowServer[0].Host, "192.10.10.10")
	assert.Equal(netflow.Server.NetflowServer[0].Port, 9999)
	assert.Equal(netflow.Server.NetflowServer[0].Name, "toto")
	assert.Equal(netflow.TemplateRefresh.Minutes, 15)
	assert.Equal(netflow.TemplateRefresh.Packets, 15)
	assert.Equal(netflow.ActiveTimeout, 15)

	// Test Edit netflow
	testNetflow.ActiveTimeout = 30
	testNetflow.Server.NetflowServer[0].Host = "newhost"
	testNetflow.Server.NetflowServer[0].Port = 8888
	testNetflow.Server.NetflowServer[0].Name = "newname"
	testNetflow.TemplateRefresh.Minutes = 30
	testNetflow.TemplateRefresh.Packets = 30
	testNetflow.Name = "newname"
	errEdit := testNetflow.Edit(pan)
	assert.Nil(errEdit)

	// Test Search netflow
	netflow, searchError = SearchNetflow("newname", pan)
	assert.Nil(searchError)
	assert.Equal(netflow.Name, "newname")
	assert.Equal(netflow.Server.NetflowServer[0].Host, "newhost")
	assert.Equal(netflow.Server.NetflowServer[0].Port, 8888)
	assert.Equal(netflow.Server.NetflowServer[0].Name, "newname")
	assert.Equal(netflow.TemplateRefresh.Minutes, 30)
	assert.Equal(netflow.TemplateRefresh.Packets, 30)
	assert.Equal(netflow.ActiveTimeout, 30)

	//Test CheckIfNetflowExist

	myTest := CheckIfNetflowExist("newname", pan)
	assert.Equal(myTest, true)

	myTest = CheckIfNetflowExist("tata", pan)
	assert.Equal(myTest, false)


	// Test Delete netflow
	errDelete := testNetflow.Delete(pan)
	assert.Nil(errDelete)

	// Test Search netflow
	test := CheckIfNetflowExist("newname", pan)
	assert.Equal(test, false)

	test = CheckIfNetflowExist("tata", pan)
	assert.Equal(test, false)

	supertest, err := SearchNetflow("netflowtest", pan)
	assert.Equal(supertest.Name, "netflowtest")
	assert.Nil(err)
	supertest.Delete(pan)

}