package api

import (
	"os"
	"testing"

	"github.com/engswee/flashpipe/internal/file"
	"github.com/engswee/flashpipe/internal/httpclnt"
	"github.com/engswee/flashpipe/internal/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APIProxySuite struct {
	suite.Suite
	serviceDetails *ServiceDetails
	exe            *httpclnt.HTTPExecuter
}

func TestAPIProxyOauth(t *testing.T) {
	suite.Run(t, &APIProxySuite{
		serviceDetails: &ServiceDetails{
			Host:              os.Getenv("FLASHPIPE_APIPORTAL_HOST"),
			OauthHost:         os.Getenv("FLASHPIPE_OAUTH_HOST"),
			OauthPath:         os.Getenv("FLASHPIPE_OAUTH_PATH"),
			OauthClientId:     os.Getenv("FLASHPIPE_APIPORTAL_OAUTH_CLIENTID"),
			OauthClientSecret: os.Getenv("FLASHPIPE_APIPORTAL_OAUTH_CLIENTSECRET"),
		},
	})
}

func (suite *APIProxySuite) SetupSuite() {
	println("========== Setting up suite - start ==========")
	suite.exe = InitHTTPExecuter(suite.serviceDetails)

	// Setup viper in case debug logs are required
	viper.SetEnvPrefix("FLASHPIPE")
	viper.AutomaticEnv()
	logger.InitConsoleLogger(viper.GetBool("debug"))

	println("========== Setting up suite - end ==========")
}

func (suite *APIProxySuite) SetupTest() {
	println("---------- Setting up test - start ----------")
	println("---------- Setting up test - end ----------")
}

func (suite *APIProxySuite) TearDownTest() {
	println("---------- Tearing down test - start ----------")
	println("---------- Tearing down test - end ----------")
}

func (suite *APIProxySuite) TearDownSuite() {
	println("========== Tearing down suite - start ==========")

	tearDownAPIProxy(suite.T(), "Northwind_V4", suite.exe)
	err := os.RemoveAll("../../output/apiproxy")
	if err != nil {
		suite.T().Logf("WARNING - Directory removal failed with error - %v", err)
	}
	println("========== Tearing down suite - end ==========")
}

func (suite *APIProxySuite) TestAPIProxy_Upload() {
	a := NewAPIProxy(suite.exe)

	err := a.Upload("../../test/testdata/apiproxy/Northwind_V4", "../../output/apiproxy/work/upload")
	if err != nil {
		suite.T().Fatalf("Upload APIProxy failed with error - %v", err)
	}
	proxyExists, err := a.Exists("Northwind_V4")
	if err != nil {
		suite.T().Fatalf("Get APIProxy failed with error %v", err)
	}
	assert.True(suite.T(), proxyExists, "APIProxy was not uploaded")

	proxies, err := a.List()
	if err != nil {
		suite.T().Fatalf("List APIProxies failed with error - %v", err)
	}
	assert.GreaterOrEqual(suite.T(), len(proxies), 1, "Expected number of APIProxies >= 1")
}

func (suite *APIProxySuite) TestAPIProxy_Download() {
	a := NewAPIProxy(suite.exe)

	err := a.Download("HelloWorldAPI", "../../output/apiproxy/work/download")
	if err != nil {
		suite.T().Fatalf("Download APIProxy failed with error - %v", err)
	}

	assert.True(suite.T(), file.Exists("../../output/apiproxy/work/download/HelloWorldAPI"), "APIProxy was not downloaded")
}

func tearDownAPIProxy(t *testing.T, id string, exe *httpclnt.HTTPExecuter) {
	a := NewAPIProxy(exe)

	proxyExists, err := a.Exists(id)
	if err != nil {
		t.Logf("WARNING - Exists failed with error - %v", err)
	}
	if proxyExists {
		err = a.Delete(id)
		if err != nil {
			t.Logf("WARNING - Delete failed with error - %v", err)
		}
	}
}
