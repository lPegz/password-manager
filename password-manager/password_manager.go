package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/lpegoraro/password-manager/remote"
)

type ConfigArgument struct {
	singleLetter, name, description string
}

var ListOfCommands = "\n\thelp | -h: Prints this message\n" +
	"\tversion | -v: Print the version of the app\n" +
	"\tget | -g {DESCRIPTION} {USERNAME} {OPTIONS}: Copy the password to the clipboard, for more information use `password_manager get help\n" +
	"\tadd | -a {DESCRIPTION} {USERNAME} {OPTIONS}: Add a new password entry, for more information use `password_manager add help\n" +
	"\tconfig | -c {METHOD} {SEED} {FACTOR} {STORAGE_TYPE} {LENGTH}: Configure encryption or password generation method\n" +
	"\t | \t \"Method\": Type of password, please choose from the following {uuid | hash | custom}\n" +
	"\t | \t \"Seed\": Any passfrase you would like\n" +
	"\t | \t \"Factor\": Given the Method uuid, you can choose between 4 and 5\n" +
	"\t | \t \"Password Length\": Given the Method hash, you define the lenght of the hash\n" +
	"\t | \t \"Storage Type\": Either {NOT_ENCRYPTED_FILE | IMMUDB }\n"

var HELP_COMMAND = ConfigArgument{
	singleLetter: "-h",
	name:         "help",
	description: "Usage: `password_manager {COMMANDS} {OPTIONS}`" +
		"\n The command list is the below" +
		ListOfCommands,
}

var ABOUT = ConfigArgument{
	singleLetter: "-v",
	name:         "version",
	description:  "Password Manager in Go version 0.1.1",
}

func main() {
	arguments := os.Args[1:]
	parseArgs(arguments)
}

var GET_COMMAND = ConfigArgument{
	singleLetter: "-g",
	name:         "get",
	description: "Copy the password to the clipboard\n" +
		"Options:\n" +
		"\thelp: Print this help\n" +
		"\tverbose: Print the password in the console, not only copy to clipboard\n",
}

var ADD_COMMAND = ConfigArgument{
	singleLetter: "-a",
	name:         "add",
	description: "Add a new password entry\n" +
		":\n" +
		"\t",
}

var SERVE_COMMAND = ConfigArgument{
	singleLetter: "-s",
	name:         "serve",
	description: "Serves a new tenant for managing your passwords\n" +
		"Options:\n" +
		"\tserve | -s {TENANT} {PUBLIC_KEY}: Returns connection instructions, to connect and serve,\n" +
		"\t | \tthe server will remain up for 5 minutes\n",
}

var CONFIG_COMMAND = ConfigArgument{
	singleLetter: "-c",
	name:         "config",
	description: "Configure the password generation\n" +
		"Options:\n" +
		"\thelp: Print this help\n" +
		"\tconfig | -c {METHOD} {SEED} {FACTOR} {STORAGE_TYPE}: Configure encryption or password generation method\n" +
		"\t | \t \"Method\": Type of password, please choose from the following {uuid | cert | custom }\n" +
		"\t | \t \"Seed\": Any passfrase you would like\n" +
		"\t | \t \"Factor\": Given the Method uuid, you can choose between 4 and 5\n" +
		"\t | \t \t   Given the Method cert you can choose the algorithym for the password creation\n" +
		"\t | \t \"Storage Type\": Only supporting \"NOT_ENCRYPTED_FILE\" storage at the moment, you can choose \n" +
		"\t | \t \t   You can choose output also, but you will need to manually configure in the settings since this \n" +
		"\t | \t \tis a development feature only.\n",
}

func parseArgs(arguments []string) error {
	command := ""
	for index := 0; index < len(arguments); index++ {
		value := arguments[index]
		if checkIfCommand(value, HELP_COMMAND) {
			command = "help"
			log.Println(HELP_COMMAND.description)
		} else if checkIfCommand(value, ABOUT) {
			command = "about"
			log.Println(ABOUT.description)
		} else if checkIfCommand(value, GET_COMMAND) {
			command = "get"
			handleGet(arguments)
		} else if checkIfCommand(value, ADD_COMMAND) {
			command = "get"
			handleAdd(arguments)
		} else if checkIfCommand(value, CONFIG_COMMAND) {
			command = "config"
			handleConfig(arguments)
		} else if checkIfCommand(value, SERVE_COMMAND) {
			command = "serve"
			handleServe(arguments)
		}
	}
	if len(arguments) == 0 || command == "" {
		log.Println(ABOUT.description)
		log.Println(HELP_COMMAND.description)
		return errors.New("No valid argument found")
	}
	return nil
}
func handleConfig(arguments []string) {
	method := arguments[1]
	seed := arguments[2]
	factor, err := strconv.ParseInt(arguments[3], 10, 8)
	if err != nil {
		log.Println(err)
	}
	storageType := arguments[4]
	CreateConfigFile(method, seed, int32(factor), storageType)
}

func handleAdd(arguments []string) {
	description := arguments[1]
	username := arguments[2]
	AddPassword(description, username, PasswordConfiguration{})
}

func handleGet(arguments []string) {
	description := arguments[1]
	username := arguments[2]
	savedPassword := GetPassword(description, username)
	if savedPassword == "" {
		log.Println("Failed to fetch password")
	}
	log.Println(savedPassword)
}

func handleServe(arguments []string) {
	tenantName := arguments[1]
	publicKey := arguments[2]
	remote.ServeNewTenant(tenantName, publicKey)
}

func AddPassword(description, username string, overrideConfig PasswordConfiguration) string {
	var configuration = PasswordConfiguration{}
	if (PasswordConfiguration{}) == overrideConfig {
		configuration = GetCurrentConfiguration()
	} else {
		configuration = overrideConfig
	}
	passwordGenerated := GeneratePassword(configuration)
	save(description, username, configuration, passwordGenerated)
	return passwordGenerated
}

func GetPassword(description string, username string) string {
	configuration := GetCurrentConfiguration()
	gotPassword := get(description, username, configuration)
	return gotPassword
}

func checkIfCommand(value string, command ConfigArgument) bool {
	return value == command.singleLetter || value == "--"+command.name || value == command.name
}
