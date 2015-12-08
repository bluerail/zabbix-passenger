package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"gopkg.in/alecthomas/kingpin.v2"
	"launchpad.net/xmlpath"
	"log"
	"os"
	"os/exec"
)

const (
	VERSION = "1.0.2"
)

func read_xml() *xmlpath.Node {
	path, err := exec.LookPath("passenger-status")
	if err != nil {
		// passenger-status not found in path
		if _, err := os.Stat("/usr/local/rvm/wrappers/default/passenger-status"); err == nil {
			// default rvm wrapper exists so use that!
			path = "/usr/local/rvm/wrappers/default/passenger-status"
		}
	}

	cmd := exec.Command(path, "--show=xml")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// Stuff to handle the iso-8859-1 xml encoding
	// http://stackoverflow.com/a/32224438/606167
	decoder := xml.NewDecoder(stdout)
	decoder.CharsetReader = charset.NewReaderLabel

	xmlData, err := xmlpath.ParseDecoder(decoder)
	if err != nil {
		log.Fatal(err)
	}

	// Check version
	version_path := xmlpath.MustCompile("/info/@version")
	if version, ok := version_path.String(xmlData); !ok || (version != "3" && version != "2") {
		log.Fatal("Unsupported Passenger version (xml version ", version, ")")
	}

	return xmlData
}

func print_simple_selector(pathString string) {
	path := xmlpath.MustCompile(pathString)

	if value, ok := path.String(read_xml()); ok {
		fmt.Println(value)
	}
}

func print_app_groups_json() {
	path := xmlpath.MustCompile("//supergroup/name")

	app_iter := path.Iter(read_xml())

	var entries []map[string]string

	for app_iter.Next() {
		entries = append(entries, map[string]string{"{#NAME}": app_iter.Node().String()})
	}

	data := map[string][]map[string]string{"data": entries}

	json, _ := json.Marshal(data)
	fmt.Println(string(json))
}

var (
	app     = kingpin.New("zabbix-passenger", "A utility to parse passenger-status output for usage with Zabbix")
	appPath = app.Flag("app", "Full path to application (leave out for global value)").String()

	appGroupsJson = app.Command("app-groups-json", "Get list of application groups in JSON format (for LLD)")
	queue         = app.Command("queue", "Get number of requests in queue, optionally specify app with --app")
	capacityUsed  = app.Command("capacity-used", "Get global capacity used, optionally specify app with --app")
)

func main() {
	app.Version(VERSION)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case appGroupsJson.FullCommand():
		print_app_groups_json()
	case queue.FullCommand():
		if *appPath != "" {
			print_simple_selector(fmt.Sprintf("//supergroup[name='%v']/group/get_wait_list_size", *appPath))
		} else {
			print_simple_selector("//info/get_wait_list_size")
		}
	case capacityUsed.FullCommand():
		if *appPath != "" {
			print_simple_selector(fmt.Sprintf("//supergroup[name='%v']/capacity_used", *appPath))
		} else {
			print_simple_selector("//info/capacity_used")
		}
	}
}
