package reading

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aewens/anda/pkg/core"
)

func Config(path string) (*core.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	raw := make(map[string]interface{})
	err = json.NewDecoder(file).Decode(&raw)
	if err != nil {
		return nil, err
	}

	rawApiPort, ok := raw["port"]
	if !ok {
		panic(fmt.Sprintf("Missing 'port' from config file: %s", path))
	}

	apiPort, valid := rawApiPort.(float64)
	if !valid {
		panic(fmt.Sprintf("Value of 'port' is not a number: %v", apiPort))
	}

	rawDatabase, ok := raw["database"]
	if !ok {
		panic(fmt.Sprintf("Missing 'database' from config file: %s", path))
	}

	database, valid := rawDatabase.(map[string]interface{})
	if !valid {
		panic(fmt.Sprintf("Value of 'port' is not a dictionary: %v",
			rawDatabase))
	}

	rawHost, ok := database["host"]
	if !ok {
		panic(fmt.Sprintf("Missing 'host' from database in config file: %s",
			path))
	}

	host, valid := rawHost.(string)
	if !valid {
		panic(fmt.Sprintf("Value of 'host' is not a string: %v", rawHost))
	}

	rawPort, ok := database["port"]
	if !ok {
		panic(fmt.Sprintf("Missing 'port' from database in config file: %s",
			path))
	}

	port, valid := rawPort.(string)
	if !valid {
		panic(fmt.Sprintf("Value of 'port' is not a string: %v", rawPort))
	}

	rawUser, ok := database["user"]
	if !ok {
		panic(fmt.Sprintf("Missing 'user' from database in config file: %s",
			path))
	}

	user, valid := rawUser.(string)
	if !valid {
		panic(fmt.Sprintf("Value of 'user' is not a string: %v", rawUser))
	}

	rawPassword, ok := database["password"]
	if !ok {
		panic(fmt.Sprintf("Missing 'password' from database in config file: %s",
			path))
	}

	password, valid := rawPassword.(string)
	if !valid {
		panic(fmt.Sprintf("Value of 'password' is not a string: %v",
			rawPassword))
	}

	rawName, ok := database["name"]
	if !ok {
		panic(fmt.Sprintf("Missing 'name' from database in config file: %s",
			path))
	}

	name, valid := rawName.(string)
	if !valid {
		panic(fmt.Sprintf("Value of 'name' is not a string: %v", rawUser))
	}

	config := &core.Config{
		ApiPort: int(apiPort),
		DBHost:  host,
		DBPort:  port,
		DBUser:  user,
		DBPswd:  password,
		DBName:  name,
	}

	return config, err
}
