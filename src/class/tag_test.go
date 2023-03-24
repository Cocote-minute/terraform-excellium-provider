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

func TestTag(t *testing.T) {
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
	
	// Tag tests
	testTag := CreateTagClass("test_unit", "color1", "test_unit")
	assert.Equal(testTag.Name,"test_unit", "Name should be test_unit and is: " + testTag.Name)
	assert.Equal(testTag.Color, "color1", "Color should be color1 and is: " + testTag.Color)
	assert.Equal(testTag.Comments, "test_unit", "Comments should be test_unit and is: " + testTag.Comments)

	// Test Add tag
	errAdd := testTag.Add(pan)
	assert.Nil(errAdd)

	// Test Search tag
	tag, searchError := SearchTag("test_unit", pan)
	assert.Nil(searchError)
	assert.Equal(tag.Name, "test_unit", "Name should be test_unit and is: " + tag.Name)
	assert.Equal(tag.Color, "color1", "Color should be color1 and is: " + tag.Color)
	assert.Equal(tag.Comments, "test_unit", "Comments should be test_unit and is: " + tag.Comments)
	
	// Test Edit tag
	testTag.Color = "color2"
	errEdit := testTag.Edit(pan)
	assert.Nil(errEdit)

	// Test Search tag
	tag, searchError = SearchTag("test_unit", pan)
	assert.Nil(searchError)
	assert.Equal(tag.Name, "test_unit", "Name should be test_unit and is: " + tag.Name)
	assert.Equal(tag.Color, "color2", "Color should be color2 and is: " + tag.Color)
	assert.Equal(tag.Comments, "test_unit", "Comments should be test_unit and is: " + tag.Comments)

	testTag.Name = "test_unit2"
	errEdit = testTag.Edit(pan)
	assert.Nil(errEdit)

	// Test Search tag
	tag, searchError = SearchTag("test_unit2", pan)
	assert.Nil(searchError)
	assert.Equal(tag.Name, "test_unit2", "Name should be test_unit2 and is: " + tag.Name)
	assert.Equal(tag.Color, "color2", "Color should be color2 and is: " + tag.Color)
	assert.Equal(tag.Comments, "test_unit", "Comments should be test_unit and is: " + tag.Comments)

	// Test CheckIfTagExist
	exsist := CheckIfTagExist("test_unit2", pan)
	assert.True(exsist, "Tag should exist")

	exsist = CheckIfTagExist("test_unit", pan)
	assert.False(exsist, "Tag should not exist")

	// Test Delete tag
	errDelete := testTag.Delete(pan)
	assert.Nil(errDelete)

	// Test check if tag exist
	exsist = CheckIfTagExist("test_unit2", pan)
	assert.False(exsist, "Tag should not exist")
}