package odata

import (
	"github.com/engswee/flashpipe/httpclnt"
	"github.com/engswee/flashpipe/logger"
	"github.com/engswee/flashpipe/str"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type PackageSuite struct {
	suite.Suite
	serviceDetails *ServiceDetails
	exe            *httpclnt.HTTPExecuter
}

func TestPackageBasicAuth(t *testing.T) {
	suite.Run(t, &PackageSuite{
		serviceDetails: &ServiceDetails{
			Host:     os.Getenv("HOST_TMN"),
			Userid:   os.Getenv("BASIC_USERID"),
			Password: os.Getenv("BASIC_PASSWORD"),
		},
	})
}

func TestPackageOauth(t *testing.T) {
	suite.Run(t, &PackageSuite{
		serviceDetails: &ServiceDetails{
			Host:              os.Getenv("HOST_TMN"),
			OauthHost:         os.Getenv("HOST_OAUTH"),
			OauthPath:         os.Getenv("HOST_OAUTH_PATH"),
			OauthClientId:     os.Getenv("OAUTH_CLIENTID"),
			OauthClientSecret: os.Getenv("OAUTH_CLIENTSECRET"),
		},
	})
}

func (suite *PackageSuite) SetupSuite() {
	println("========== Setting up suite ==========")
	suite.exe = InitHTTPExecuter(suite.serviceDetails)

	// Setup viper in case debug logs are required
	viper.SetEnvPrefix("FLASHPIPE")
	viper.AutomaticEnv()
	logger.InitConsoleLogger(viper.GetBool("debug"))

	setupPackage(suite.T(), "FlashPipeIntegrationTest", suite.exe)

	setupArtifact(suite.T(), "IFlow1", "FlashPipeIntegrationTest", "../testdata/artifacts/setup/IFlow1", "Integration", suite.exe)
}

func (suite *PackageSuite) SetupTest() {
	println("---------- Setting up test ----------")
}

func (suite *PackageSuite) TearDownTest() {
	println("---------- Tearing down test ----------")
}

func (suite *PackageSuite) TearDownSuite() {
	println("========== Tearing down suite ==========")

	tearDownPackage(suite.T(), "FlashPipeIntegrationTest", suite.exe)
	tearDownPackage(suite.T(), "FlashPipeIntegrationTestCreate", suite.exe)
}

func (suite *PackageSuite) TestIntegrationPackage_CreateUpdateDelete() {
	const packageId = "FlashPipeIntegrationTestCreate"
	ip := NewIntegrationPackage(suite.exe)

	jsonData := new(PackageSingleData)
	jsonData.Root.Id = packageId
	jsonData.Root.Name = "FlashPipe Integration Test Create"
	jsonData.Root.ShortText = "FlashPipe Integration Test Create"
	jsonData.Root.Mode = "EDIT_ALLOWED"
	// Create
	err := ip.Create(jsonData)
	if err != nil {
		suite.T().Fatalf("Create failed with error - %v", err)
	}

	// Update
	jsonData.Root.Name = "FlashPipe Integration Test Update"
	jsonData.Root.Name = "FlashPipe Integration Test Update"
	jsonData.Root.Mode = "EDIT_ALLOWED"
	err = ip.Update(jsonData)
	if err != nil {
		suite.T().Fatalf("Update failed with error - %v", err)
	}

	// Get list
	packagesList, err := ip.GetPackagesList()
	if err != nil {
		suite.T().Fatalf("GetPackagesList failed with error - %v", err)
	}
	assert.Truef(suite.T(), str.Contains(packageId, packagesList), "%v found in packagesList", packageId)

	// Check not read only
	_, readOnly, _, err := ip.Get(packageId)
	if err != nil {
		suite.T().Fatalf("IsReadOnly failed with error - %v", err)
	}
	assert.Falsef(suite.T(), readOnly, "%v is not read only", packageId)

	// Delete
	err = ip.Delete(packageId)
	if err != nil {
		suite.T().Fatalf("Delete failed with error - %v", err)
	}
}

func (suite *PackageSuite) TestIntegrationPackage_GetArtifacts() {
	ip := NewIntegrationPackage(suite.exe)

	artifacts, err := ip.GetAllArtifacts("FlashPipeIntegrationTest")
	if err != nil {
		suite.T().Fatalf("GetAllArtifacts failed with error - %v", err)
	}
	artifact := FindArtifactById("IFlow1", artifacts)
	assert.NotNil(suite.T(), artifact, "IFlow1 found")
}

func setupPackage(t *testing.T, packageId string, exe *httpclnt.HTTPExecuter) {
	ip := NewIntegrationPackage(exe)

	_, _, packageExists, err := ip.Get(packageId)
	if err != nil {
		t.Fatalf("Exists failed with error - %v", err)
	}
	if !packageExists {
		requestBody := new(PackageSingleData)
		requestBody.Root.Id = packageId
		requestBody.Root.Name = packageId
		requestBody.Root.ShortText = packageId

		err = ip.Create(requestBody)
		if err != nil {
			t.Fatalf("Create failed with error - %v", err)
		}
	}
}

func tearDownPackage(t *testing.T, packageId string, exe *httpclnt.HTTPExecuter) {
	ip := NewIntegrationPackage(exe)

	_, _, packageExists, err := ip.Get(packageId)
	if err != nil {
		t.Fatalf("Exists failed with error - %v", err)
	}
	if packageExists {
		err = ip.Delete(packageId)
		if err != nil {
			t.Fatalf("Delete failed with error - %v", err)
		}
	}
}