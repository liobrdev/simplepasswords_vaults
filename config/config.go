package config

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

type AppConfig struct {
	API_GATEWAY_HOST		string
	API_GATEWAY_PORT		string
	ENVIRONMENT					string
	PASSWORD_HEADER_KEY	string
	VAULTS_ACCESS_TOKEN	string
	VAULTS_DB_HOST			string
	VAULTS_DB_NAME			string
	VAULTS_DB_PASSWORD	string
	VAULTS_DB_PORT			string
	VAULTS_DB_USER			string
	VAULTS_HOST					string
	VAULTS_PORT					string
	GO_TESTING_CONTEXT	*testing.T
}

type envAbsPaths struct {
	API_GATEWAY_HOST		string
	API_GATEWAY_PORT		string
	ENVIRONMENT					string
	PASSWORD_HEADER_KEY	string
	VAULTS_ACCESS_TOKEN	string
	VAULTS_DB_HOST			string
	VAULTS_DB_NAME			string
	VAULTS_DB_PASSWORD	string
	VAULTS_DB_PORT			string
	VAULTS_DB_USER			string
	VAULTS_HOST					string
	VAULTS_PORT					string
}

func scanFileFirstLineToConf(file *os.File, confElem *reflect.Value, path, fieldName string) {
	scanner := bufio.NewScanner(file)
	scanner.Scan()

	if contents := scanner.Text(); scanner.Err() != nil {
		log.Fatalf(
			"Error reading contents of '%s' from environment variable %s:\n%s",
			path, fieldName, scanner.Err(),
		)
	} else if contents == "" {
		log.Fatalf("Empty contents of '%s' from environment variable %s", path, fieldName)
	} else {
		confElem.FieldByName(fieldName).SetString(contents)
	}
}

func loadFileContentsFromPathsToConf(
	conf *AppConfig, pathsType *reflect.Type, pathsValue *reflect.Value, fieldCount int,
) {
	confElem := reflect.ValueOf(conf).Elem()

	for i := 0; i < fieldCount; i++ {
		fieldName := (*pathsType).Field(i).Name
		path := pathsValue.Field(i).Interface().(string)

		if path == "" {
			log.Fatal("Missing or empty environment variable: ", fieldName)
		}

		file, err := os.Open(path)

		if err != nil {
			log.Fatalf(
				"Error opening '%s' from environment variable %s:\n%s", path, fieldName, err,
			)
		}

		defer file.Close()

		scanFileFirstLineToConf(file, &confElem, path, fieldName)
	}
}

func LoadConfigFromEnv(conf *AppConfig) (err error) {
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if err.Error() == "open ./.env: no such file or directory" {
			err = nil
		} else {
			return
		}
	}

	paths := envAbsPaths{}
	pathsValue := reflect.ValueOf(paths)
	pathsType := pathsValue.Type()
	fieldCount := pathsValue.NumField()

	for i, key, val := 0, "", ""; i < fieldCount; i++ {
		key = pathsType.Field(i).Name

		if val = os.Getenv(key); val != "" {
			viper.Set(key, val)
		}
	}

	if err = viper.Unmarshal(&paths); err != nil {
		return
	} else {
		pathsValue = reflect.ValueOf(paths)
	}

	loadFileContentsFromPathsToConf(conf, &pathsType, &pathsValue, fieldCount)

	return
}
