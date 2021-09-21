package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/twpayne/go-kml"
)

type GeoPlugin struct {
	Geoplugin_request                string
	Geoplugin_status                 int
	Geoplugin_delay                  string
	Geoplugin_credit                 string
	Geoplugin_city                   string
	Geoplugin_region                 string
	Geoplugin_regionCode             string
	Geoplugin_regionName             string
	Geoplugin_areaCode               string
	Geoplugin_dmaCode                string
	Geoplugin_countryCode            string
	Geoplugin_countryName            string
	Goplugin_inEU                    int
	Geoplugin_euVATrate              bool
	Geoplugin_continentCode          string
	Geoplugin_continentName          string
	Geoplugin_latitude               string
	Geoplugin_longitude              string
	Geoplugin_locationAccuracyRadius string
	Geoplugin_timezone               string
	Geoplugin_currencyCode           string
	Geoplugin_currencySymbol         string
	Geoplugin_currencySymbol_UTF8    string
	Geoplugin_currencyConverter      int
}

type IpHub struct {
	Ip          string
	CountryCode string
	CountryName string
	Asn         int
	Isp         string
	Block       int
	Hostname    string
}

func main() {
	ips := readIPsFromFile("ipFile.txt")
	fmt.Println(ips)
	concatStats := ""
	concatGeo := ""
	var strStatsTitle string = "IP Address, Latitude, Longitude, City, Region, Country, Timezone, ISP, Blocked\n"
	concatStats += strStatsTitle
	c := make(chan []byte)
	for _, ip := range ips {
		fmt.Println(ip)
		go handleGETRequests(ip, c)
		geoPluginBody := <-c
		ipHubBody := <-c
		var geo GeoPlugin
		var iphub IpHub
		json.Unmarshal(geoPluginBody, &geo)
		json.Unmarshal(ipHubBody, &iphub)
		strGeoTitle := geo.Geoplugin_request + " " + geo.Geoplugin_city
		strGeo := strGeoTitle + "," + geo.Geoplugin_latitude + "," + geo.Geoplugin_longitude + "\n"
		strStats := geo.Geoplugin_request + "," + geo.Geoplugin_latitude + "," + geo.Geoplugin_longitude + "," +
			geo.Geoplugin_city + "," + geo.Geoplugin_region + "," + geo.Geoplugin_countryName + "," +
			geo.Geoplugin_timezone + "," + iphub.Isp + "," + strconv.Itoa(iphub.Block) + "\n"
		concatStats += strStats
		concatGeo += strGeo
	}

	saveToFile("geo-coords.csv", concatGeo)
	saveToFile("geo-stats.csv", concatStats)
	createKMLFile()

}

func handleGETRequests(ip string, c chan []byte) {
	resp := requestGeoPlugin(ip)
	body, _ := ioutil.ReadAll(resp.Body)
	resp2 := getVPNStatus(ip)
	body2, _ := ioutil.ReadAll(resp2.Body)
	c <- body
	c <- body2
}

func requestGeoPlugin(ip string) *http.Response {
	resp, err := http.Get("http://www.geoplugin.net/json.gp?ip=" + ip)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func getVPNStatus(ip string) *http.Response {
	strLink := "http://v2.api.iphub.info/ip/" + ip
	client := &http.Client{}
	req, _ := http.NewRequest("GET", strLink, nil)
	strKey := readIpHubKeyFromFile()
	req.Header.Set("X-Key", strKey)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func readIPsFromFile(filename string) []string {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	s := strings.Split(string(bs), "\n")
	return s
}

func readIpHubKeyFromFile() string {
	bs, err := ioutil.ReadFile("iphub-info.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(bs))
}

func saveToFile(filename string, d string) error {
	return ioutil.WriteFile(filename, []byte(d), 0666)
}

func createKMLFile() {
	var arPlacemarks []kml.Element
	bs, _ := ioutil.ReadFile("geo-coords.csv")
	s := strings.Split(string(bs), "\n")

	for i := 0; i < len(s)-1; i++ {
		line := strings.Split(s[i], ",")
		floatLong, _ := strconv.ParseFloat(line[2], 16) // must be switched to work
		floatLat, _ := strconv.ParseFloat(line[1], 16)

		k := kml.Placemark(
			kml.Name(line[0]),
			kml.Point(
				kml.Coordinates(kml.Coordinate{Lon: floatLong, Lat: floatLat}),
			),
		)
		arPlacemarks = append(arPlacemarks, k)

	}

	dt := kml.KML(
		kml.Document(
			arPlacemarks...,
		),
	)
	currentDirectory, _ := os.Getwd()
	fname := filepath.Join(currentDirectory, "/", "locations.kml")
	temp, _ := os.Create(fname)

	if err := dt.WriteIndent(temp, "", " "); err != nil {
		log.Fatal(err)
	}
	log.Println("Success! locations.xml is located at: " + fname)
}
