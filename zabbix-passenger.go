package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"golang.org/x/net/html/charset"
	"launchpad.net/xmlpath"
	"log"
	"os"
)

func read_xml() *xmlpath.Node {
	// TODO Execute passenger-status
	reader, err := os.Open("ppm.xml")
	if err != nil {
		log.Fatal(err)
	}

	// Stuff to handle the iso-8859-1 xml encoding
	// http://stackoverflow.com/a/32224438/606167
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	xmlData, err := xmlpath.ParseDecoder(decoder)
	if err != nil {
		log.Fatal(err)
	}

	// Check version
	version_path := xmlpath.MustCompile("/info/@version")
	if version, ok := version_path.String(xmlData); !ok || version != "3" {
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

	fmt.Println("{\"data\": [")
	for app_iter.Next() {
		fmt.Printf("{\"{#NAME}\": \"%v\"},\n", app_iter.Node().String())
	}
	fmt.Println("]}")
}

func main() {
	appGroupsJson := flag.Bool("app-groups-json", false, "Show application groups in JSON format")
	globalQueue := flag.Bool("global-queue", false, "Print number of requests in global queue")
	appQueue := flag.String("app-queue", "", "Print number of requests in specified app queue")
	globalCapacityUsed := flag.Bool("global-capacity-used", false, "Print global capacity used")
	appCapacityUsed := flag.String("app-capacity-used", "", "Print specified app capacity used")

	flag.Parse()

	if *appGroupsJson {
		print_app_groups_json()
	} else if *globalQueue {
		print_simple_selector("//info/get_wait_list_size")
	} else if *appQueue != "" {
		print_simple_selector(fmt.Sprintf("//supergroup[name='%v']/get_wait_list_size", *appQueue))
	} else if *globalCapacityUsed {
		print_simple_selector("//info/capacity_used")
	} else if *appCapacityUsed != "" {
		print_simple_selector(fmt.Sprintf("//supergroup[name='%v']/capacity_used", *appCapacityUsed))
	} else {
		flag.Usage()
	}
}
