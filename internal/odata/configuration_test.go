package odata

import (
	"github.com/engswee/flashpipe/internal/httpclnt"
	"github.com/engswee/flashpipe/internal/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ConfigurationSuite struct {
	suite.Suite
	serviceDetails *ServiceDetails
	exe            *httpclnt.HTTPExecuter
}

func TestConfigurationBasicAuth(t *testing.T) {
	suite.Run(t, &ConfigurationSuite{
		serviceDetails: &ServiceDetails{
			Host:     os.Getenv("HOST_TMN"),
			Userid:   os.Getenv("BASIC_USERID"),
			Password: os.Getenv("BASIC_PASSWORD"),
		},
	})
}

func TestConfigurationOauth(t *testing.T) {
	suite.Run(t, &ConfigurationSuite{
		serviceDetails: &ServiceDetails{
			Host:              os.Getenv("HOST_TMN"),
			OauthHost:         os.Getenv("HOST_OAUTH"),
			OauthPath:         os.Getenv("HOST_OAUTH_PATH"),
			OauthClientId:     os.Getenv("OAUTH_CLIENTID"),
			OauthClientSecret: os.Getenv("OAUTH_CLIENTSECRET"),
		},
	})
}

func (suite *ConfigurationSuite) SetupSuite() {
	println("========== Setting up suite ==========")
	suite.exe = InitHTTPExecuter(suite.serviceDetails)

	// Setup viper in case debug logs are required
	viper.SetEnvPrefix("FLASHPIPE")
	viper.AutomaticEnv()
	logger.InitConsoleLogger(viper.GetBool("debug"))

	setupPackage(suite.T(), "FlashPipeIntegrationTest", suite.exe)

	setupArtifact(suite.T(), "Integration_Test_IFlow", "FlashPipeIntegrationTest", "../../test/testdata/artifacts/update/Integration_Test_IFlow", "Integration", suite.exe)
}

func (suite *ConfigurationSuite) SetupTest() {
	println("---------- Setting up test ----------")
}

func (suite *ConfigurationSuite) TearDownTest() {
	println("---------- Tearing down test ----------")
}

func (suite *ConfigurationSuite) TearDownSuite() {
	println("========== Tearing down suite ==========")

	tearDownPackage(suite.T(), "FlashPipeIntegrationTest", suite.exe)
}

func (suite *ConfigurationSuite) TestConfiguration_Get() {
	c := NewConfiguration(suite.exe)

	parametersData, err := c.Get("Integration_Test_IFlow", "active")
	if err != nil {
		suite.T().Fatalf("Get failed with error - %v", err)
	}
	parameter := FindParameterByKey("Sender Endpoint", parametersData.Root.Results)
	assert.Equal(suite.T(), "/flow", parameter.ParameterValue, "Parameter Sender Endpoint should have value /flow")
}

func (suite *ConfigurationSuite) TestConfiguration_Update() {
	c := NewConfiguration(suite.exe)

	err := c.Update("Integration_Test_IFlow", "active", "Sender Endpoint", "/flow_update")
	if err != nil {
		suite.T().Fatalf("Update failed with error - %v", err)
	}
	parametersData, err := c.Get("Integration_Test_IFlow", "active")
	if err != nil {
		suite.T().Fatalf("Get failed with error - %v", err)
	}
	parameter := FindParameterByKey("Sender Endpoint", parametersData.Root.Results)
	assert.Equal(suite.T(), "/flow_update", parameter.ParameterValue, "Parameter Sender Endpoint should have value /flow_update after update")
}