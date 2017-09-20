package main

import (
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/cf/flags"
)


type MultiCUPSPlugin struct{}

func (c *MultiCUPSPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command multi-cups-plugin
	fc, err := parseArguments(args)

	if err != nil {
		exit1(err.Error())
	}

	if args[0] == "multi-cups-plugin" {
		fmt.Println("Running the multi-cups-plugin")
		fmt.Println(fc.String("path"))
		loadCUPS(fc.String("path"), cliConnection)
	}
}

func (c *MultiCUPSPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "MyMultiCUPSPlugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "multi-cups-plugin",
				Alias: "mcups",
				HelpText: "A command to help you create multip cups services from a bigger json file",
				UsageDetails: plugin.Usage{
					Usage: "multi-cups-plugin\n   cf multi-cups-plugin",
				},
			},
		},
	}
}

func parseArguments(args []string) (flags.FlagContext, error) {
	fc := flags.New()
	fc.NewStringFlag("path", "p", "path to cups json")
	err := fc.Parse(args...)

	return fc, err
}

func main() {
	plugin.Start(new(MultiCUPSPlugin))
}


func exit1(err string) {
	fmt.Println("FAILED\n" + err)
	os.Exit(1)
}

type CredEntry struct {
		Credentials *json.RawMessage `json:"credentials"`
		Name   string           `json:"name"`
}

func loadCUPS(file string, cliConnection plugin.CliConnection) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("unmarshal error")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var f []CredEntry
	err = json.Unmarshal(raw, &f)
	if err != nil {
		fmt.Println("READ ERROR: Is your json valid?")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, credEntry := range f {
		b, err := credEntry.Credentials.MarshalJSON()
		if err != nil {
			fmt.Println("marshal error")
			fmt.Println(err.Error())
			continue
		}
		//Check if service already exists
		service, err := cliConnection.GetService(credEntry.Name)
		fmt.Println(service)
		if err == nil {
			fmt.Println("Updating Service")
			cliConnection.CliCommand("update-user-provided-service", credEntry.Name, "-p", string(b))
		} else {
			fmt.Println("Create New Service")
			cliConnection.CliCommand("create-user-provided-service", credEntry.Name, "-p", string(b))
		}

	}

}

//TODO check if service exists
//If service then UUPS service
