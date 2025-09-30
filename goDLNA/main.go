package main

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/huin/goupnp"
)

type Container struct {
	ID          string `xml:"id,attr"`
	ParentID    string `xml:"parentID,attr"`
	Restricted  string `xml:"restricted,attr"`
	Searchable  string `xml:"searchable,attr"`
	ChildCount  int    `xml:"childCount,attr"`
	Title       string `xml:"title"`
	Class       string `xml:"class"`
	StorageUsed string `xml:"storageUsed"`
}

type Res struct {
	Size            string `xml:"size,attr"`
	Duration        string `xml:"duration,attr"`
	Bitrate         string `xml:"bitrate,attr"`
	SampleFrequency string `xml:"sampleFrequency,attr"`
	AudioChannels   string `xml:"nrAudioChannels,attr"`
	Resolution      string `xml:"resolution,attr"`
	ProtocolInfo    string `xml:"protocolInfo,attr"`
	URL             string `xml:",chardata"`
}

type Item struct {
	ID            string `xml:"id,attr"`
	ParentID      string `xml:"parentID,attr"`
	Restricted    string `xml:"restricted,attr"`
	Title         string `xml:"http://purl.org/dc/elements/1.1/ title"`
	Class         string `xml:"urn:schemas-upnp-org:metadata-1-0/upnp/ class"`
	Date          string `xml:"http://purl.org/dc/elements/1.1/ date"`
	EpisodeSeason string `xml:"urn:schemas-upnp-org:metadata-1-0/upnp/ episodeSeason"`
	EpisodeNumber string `xml:"urn:schemas-upnp-org:metadata-1-0/upnp/ episodeNumber"`
	Res           Res    `xml:"res"`
}

type DIDLLite struct {
	XMLName    xml.Name    `xml:"DIDL-Lite"`
	Namespace  string      `xml:"xmlns,attr"`
	Containers []Container `xml:"container"`
	Items      []Item      `xml:"item"`
}

func main() {
	arg := os.Args[1]
	if arg == "--version" || arg == "-v" {
		version()
		return
	}
	if arg == "--list" || arg == "-l" {
		list()
		return
	}
	if arg == "--wake" || arg == "-w" {
		wake(os.Args[2])
		return
	}
	if arg == "--wake" || arg == "-w" {
		wake(os.Args[2])
		return
	}
	if arg == "--browse" || arg == "-b" {
		count := "2000"
		if os.Args[4] != "" {
			count = os.Args[4]
		}
		browse(os.Args[2], os.Args[3], count)
		return
	}
	if arg == "--info" || arg == "-i" {
		info()
		return
	}
	help()
}

func list() {
	list, err := goupnp.DiscoverDevices("upnp:rootdevice")
	if err != nil {
		fmt.Println(err)
	}
	for _, maybeDevice := range list {
		device, err := goupnp.DeviceByURL(maybeDevice.Location)
		if err != nil {
			continue
		}
		if strings.Contains(device.Device.DeviceType, "MediaServer") {
			fmt.Println("")
			fmt.Println(device.Device.FriendlyName)
			fmt.Println(maybeDevice.Location)
		}
	}
}

func help() {
	fmt.Println("mpvDLNA.py supports the following commands:")
	fmt.Println("-h, --help     Prints the help dialog")
	fmt.Println("-v, --version  Prints version information")
	fmt.Println("-l, --list     Takes a timeout in seconds and outputs a list of DLNA Media Servers on the network")
	fmt.Println("-b, --browse   Takes a DLNA url and the id of a DLNA element and outputs its direct children")
	fmt.Println("-i, --info     Takes a DLNA url and the id of a DLNA element and outputs its metadata")
	fmt.Println("-w, --wake     Takes a MAC address and attempts to send a wake on lan packet to it")
}

func version() {
	fmt.Println("Drop in replacement for mpvDLNA.py v 2.1.0")
}

func wake(mac string) {
	fmt.Println("Not implemented yet", mac)
}

func browse(deviceUrl string, id string, count string) {
	parsedUrl, err := url.Parse(deviceUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	device, err := goupnp.DeviceByURL(parsedUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, service := range device.Device.Services {
		if strings.Contains(service.ServiceType, "ContentDirectory") {
			client := service.NewSOAPClient()
			args := &struct {
				ObjectID       string
				BrowseFlag     string
				Filter         string
				StartingIndex  string
				RequestedCount string
				SortCriteria   string
			}{ObjectID: id, BrowseFlag: "BrowseDirectChildren", Filter: "*", StartingIndex: "0", RequestedCount: "2000", SortCriteria: ""}

			response := &struct {
				Result         string
				NumberReturned string
				TotalMatches   string
				UpdateID       string
			}{}

			err = client.PerformAction("urn:schemas-upnp-org:Service:ContentDirectory:1", "Browse", args, response)

			var didl DIDLLite
			// fmt.Println(response.Result)
			err := xml.Unmarshal([]byte(response.Result), &didl)
			if err != nil {
				fmt.Println(err)
			}

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("items:")
			for _, item := range didl.Items {
				fmt.Println("")
				fmt.Println(item.Title)
				fmt.Println(item.ID)
				fmt.Println(item.Res.URL)
			}
			fmt.Println("\x1F")
			fmt.Println("containers:")
			for _, container := range didl.Containers {
				fmt.Println("")
				fmt.Println(container.Title)
				fmt.Println(container.ID)
			}
			fmt.Println("\x1F")
		}
	}
}

func info() {
	fmt.Println("unknown")
	fmt.Println("")
	fmt.Println("No Episode Number")
	fmt.Println("No Description")
}
