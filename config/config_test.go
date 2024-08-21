package config_test

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config"
	"github.com/stretchr/testify/require"
)

const ConfigFileWithUnknownFieldTest = `
[Log]
	Environment = "development" # "production" or "development"
	Level = "info"
	Outputs = ["stderr"]
	ThisIsAnUnknownField = "unknown"
`
const ConfigFileOkTest = `
[Log]
	Environment = "development" # "production" or "development"
	Level = "info"
	Outputs = ["stderr"]
	#ThisIsAnUnknownField = "unknown"
`

const ConfigFileValidiumTranslatorTest = `
[Etherman.Validium.Translator]
			FullMatchRules = [
				{Old="http://dataavailability-003.cdk-validium-cardona-03.zkevm.polygon.private:8444", New="https://dataavailability-003-cdk-validium-cardona-03-zkevm.polygondev.tools"},
				{Old="http://dataavailability-002.cdk-validium-cardona-03.zkevm.polygon.private:8444", New="https://dataavailability-002-cdk-validium-cardona-03-zkevm.polygondev.tools"},
				{Old="http://dataavailability-001.cdk-validium-cardona-03.zkevm.polygon.private:8444", New="https://dataavailability-001-cdk-validium-cardona-01-zkevm.polygondev.tools"}
			]
`

const ConfigFileValidiumTranslatorWrongFieldsInMapTest = `
[Etherman.Validium.Translator]
			FullMatchRules = [
				{NoExpectedField="http://dataavailability-003.cdk-validium-cardona-03.zkevm.polygon.private:8444", New="https://dataavailability-003-cdk-validium-cardona-03-zkevm.polygondev.tools"},
				{Old="http://dataavailability-002.cdk-validium-cardona-03.zkevm.polygon.private:8444", New="https://dataavailability-002-cdk-validium-cardona-03-zkevm.polygondev.tools"},
				{Old="http://dataavailability-001.cdk-validium-cardona-03.zkevm.polygon.private:8444", New="https://dataavailability-001-cdk-validium-cardona-01-zkevm.polygondev.tools"}
			]
`

func TestLoadConfigUnknownFieldFails(t *testing.T) {
	fileExtension := "toml"
	_, err := config.LoadFileFromString(string(ConfigFileWithUnknownFieldTest), fileExtension)
	require.Error(t, err)
}

func TestLoadConfigUnknownFieldInMapFails(t *testing.T) {
	t.Skip("This test is not working as expected, need a deep dive into Viper library")
	fileExtension := "toml"
	_, err := config.LoadFileFromString(string(ConfigFileValidiumTranslatorWrongFieldsInMapTest), fileExtension)
	require.Error(t, err)
}

func TestLoadConfigCommentedUnknownFieldOk(t *testing.T) {
	fileExtension := "toml"
	_, err := config.LoadFileFromString(string(ConfigFileOkTest), fileExtension)
	require.NoError(t, err)
}
func TestLoadConfigValdiumTranslatorOk(t *testing.T) {
	fileExtension := "toml"
	cfg, err := config.LoadFileFromString(string(ConfigFileValidiumTranslatorTest), fileExtension)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, 3, len(cfg.Etherman.Validium.Translator.FullMatchRules))
}
