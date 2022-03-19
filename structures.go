package speedtest

import (
	"encoding/xml"
	"sort"
)

type Server struct {
	ID          int     `xml:"id,attr"`
	URL         string  `xml:"url,attr"`
	Lat         float64 `xml:"lat,attr"`
	Lon         float64 `xml:"lon,attr"`
	Name        string  `xml:"name,attr"`
	Country     string  `xml:"country,attr"`
	CountryCode string  `xml:"cc,attr"`
	Sponsor     string  `xml:"sponsor,attr"`
	
	Distance float64
}

type Servers []Server

func (s Servers) Len() int      { return len(s) }

func (s Servers) Swap(i, j int) { s[i], s[j] = s[j], s[i] }


type byDistance struct{ Servers }
func (s byDistance) Less(i, j int) bool {
	return s.Servers[i].Distance < s.Servers[j].Distance
}


type byID struct{ Servers }
func (s byID) Less(i, j int) bool {
	return s.Servers[i].ID < s.Servers[j].ID
}

func (s Servers) SortByID() {
	sort.Sort(byID{s})
}


func (s Servers) SortByDistance() {
	sort.Sort(byDistance{s})
}

type Settings struct {
	XMLName xml.Name `xml:"settings"`
	Servers Servers  `xml:"servers>server"`
}

type Config struct {
	XMLName xml.Name `xml:"settings"`
	Client  Client   `xml:"client"`
}

type Client struct {
	IPAddress string  `xml:"ip,attr"`
	Lat       float64 `xml:"lat,attr"`
	Lon       float64 `xml:"lon,attr"`
	IspName   string  `xml:"isp,attr"`
}

func (settings Settings) UpdateDistances(lat float64, lon float64) {
	for i, server := range settings.Servers {
		settings.Servers[i].Distance = Distance(
			server.Lat*degToRad, server.Lon*degToRad,
			lat*degToRad, lon*degToRad)
	}
}
