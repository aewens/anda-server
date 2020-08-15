package main

import (
	"encoding/json"
	//	"encodeing/base64"
	//	"errors"
	"flag"
	"fmt"
	//	"io/ioutil"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

type Config struct {
	ApiPort int
	DBHost  string
	DBPort  string
	DBUser  string
	DBPswd  string
	DBName  string
}

type EAV struct {
	UUID  string      `json:"uuid"`
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
	Flag  int         `json:"flag"`
}

type EAVList []EAV

func cleanup() {
	fmt.Println("Closing server")
}

func HandleSigterm() {
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		cleanup()
		os.Exit(1)
	}()
}

func ParseFlags() map[string]interface{} {
	configFlag := flag.String("config", "", "Path to config file")
	flag.Parse()

	if len(*configFlag) == 0 {
		log.Println("ERROR: Config flag is missing!")
		os.Exit(2)
	}

	flags := make(map[string]interface{})
	flags["config"] = *configFlag

	return flags
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()
	//buf, err := ioutil.ReadAll(file)
	//if err != nil {
	//	log.Println(err)
	//}

	raw := make(map[string]interface{})
	json.NewDecoder(file).Decode(&raw)

	rawApiPort, ok := raw["port"]
	if !ok {
		log.Fatalf("Missing 'port' from config file: %s", path)
	}

	apiPort, valid := rawApiPort.(float64)
	if !valid {
		log.Fatalf("Value of 'port' is not a number: %v", apiPort)
	}

	rawDatabase, ok := raw["database"]
	if !ok {
		log.Fatalf("Missing 'database' from config file: %s", path)
	}

	database, valid := rawDatabase.(map[string]interface{})
	if !valid {
		log.Fatalf("Value of 'port' is not a dictionary: %v", rawDatabase)
	}

	rawHost, ok := database["host"]
	if !ok {
		log.Fatalf("Missing 'host' from database in config file: %s", path)
	}

	host, valid := rawHost.(string)
	if !valid {
		log.Fatalf("Value of 'host' is not a string: %v", rawHost)
	}

	rawPort, ok := database["port"]
	if !ok {
		log.Fatalf("Missing 'port' from database in config file: %s", path)
	}

	port, valid := rawPort.(string)
	if !valid {
		log.Fatalf("Value of 'port' is not a string: %v", rawPort)
	}

	rawUser, ok := database["user"]
	if !ok {
		log.Fatalf("Missing 'user' from database in config file: %s", path)
	}

	user, valid := rawUser.(string)
	if !valid {
		log.Fatalf("Value of 'user' is not a string: %v", rawUser)
	}

	rawPassword, ok := database["password"]
	if !ok {
		log.Fatalf("Missing 'password' from database in config file: %s", path)
	}

	password, valid := rawPassword.(string)
	if !valid {
		log.Fatalf("Value of 'password' is not a string: %v", rawPassword)
	}

	rawName, ok := database["name"]
	if !ok {
		log.Fatalf("Missing 'name' from database in config file: %s", path)
	}

	name, valid := rawName.(string)
	if !valid {
		log.Fatalf("Value of 'name' is not a string: %v", rawUser)
	}

	config := &Config{
		ApiPort: int(apiPort),
		DBHost:  host,
		DBPort:  port,
		DBUser:  user,
		DBPswd:  password,
		DBName:  name,
	}

	return config, err
}

func EAVSelect(db *sql.DB) EAVList {
	var entries EAVList

	query := `
		SELECT e.uuid, a.name, vt.name, convert_from(v.value, 'utf-8'), v.flag
		FROM entity e
		INNER JOIN value v ON v.entity_id = e.id
		INNER JOIN attribute a ON a.id = v.attribute_id
		INNER JOIN value_type vt ON vt.id = a.value_type_id;
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		os.Exit(5)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			uuid  string
			attr  string
			vtype string
			value interface{}
			flag  int
		)
		err := rows.Scan(&uuid, &attr, &vtype, &value, &flag)
		if err != nil {
			log.Println(err)
			os.Exit(6)
		}
		log.Println(uuid, attr, vtype, value, flag)
		entries = append(entries, EAV{uuid, attr, vtype, value, flag})
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		os.Exit(7)
	}

	return entries
}

func Welcome(server *Server) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	}
}

func GetEntries(server *Server) Handler {

	return func(w http.ResponseWriter, r *http.Request) {
		entries := EAVSelect(server.DB)
		json.NewEncoder(w).Encode(entries)
	}
}

func StartServer(config *Config) {
	server := CreateServer(config)
	server.AddRoute("GET", "/api", Welcome)
	server.AddRoute("GET", "/api/entries", GetEntries)

	server.Start()
}

func main() {
	HandleSigterm()

	flags := ParseFlags()
	configFlag, ok := flags["config"]
	if !ok {
		log.Println("ERROR: Config flag is missing!")
		os.Exit(3)
	}

	configPath, ok := configFlag.(string)
	if !ok {
		log.Fatalf("Config flag is not a string: %v", configFlag)
	}

	config, err := ReadConfig(configPath)
	if err != nil {
		log.Println(err)
		os.Exit(4)
	}

	StartServer(config)
}
