package geomap

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"time"

	"github.com/QOLPlus/weather-map/utils"
)

const fileName = "geomap.yaml"
const repositoryUrl = "https://github.com/QOLPlus/weather-map"
const api = "https://www.weather.go.kr/weather/lifenindustry/sevice_rss.jsp"

type ParsedPair struct {
	Code string
	Name string
}
func (p ParsedPair) String() string {
	return fmt.Sprintf("%s : %s", p.Name, p.Code)
}

type GeoMap struct {
	Generated geoMapGenerated `yaml:"generated"`
	Data      []geoMapNode    `yaml:"data"`
}
type geoMapGenerated struct {
	At string `yaml:"at"`
	By string `yaml:"by"`
}
type geoMapNode struct {
	Name     string       `yaml:"name"`
	Code     string       `yaml:"code"`
	Children []geoMapNode `yaml:"children"`
}
func (gm GeoMap) export() {
	marshaled, err := yaml.Marshal(&gm)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func(){ _ = file.Close() }()

	writer := bufio.NewWriter(file)
	writtenBytes, err := writer.WriteString(string(marshaled))
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
	fmt.Printf("geomap %d bytes written!\n", writtenBytes)
}

func Process(cmd *cobra.Command, args []string) {
	data := GeoMap{
		Generated: geoMapGenerated{
			At: time.Now().Format(time.RFC3339),
			By: repositoryUrl,
		},
	}

	fmt.Printf("\nGenerating geomap started at %s !\n", data.Generated.At)

	for _, sido := range getParsedList([]ParsedPair{}) {
		sidoNode := geoMapNode{Name: sido.Name, Code: sido.Code}

		for _, gugun := range getParsedList([]ParsedPair{sido}) {
			gugunNode := geoMapNode{Name: gugun.Name, Code: gugun.Code}

			for _, dong := range getParsedList([]ParsedPair{sido, gugun}) {
				dongNode := geoMapNode{Name: dong.Name, Code: dong.Code}
				gugunNode.Children = append(gugunNode.Children, dongNode)
			}

			sidoNode.Children = append(sidoNode.Children, gugunNode)
		}

		data.Data = append(data.Data, sidoNode)
	}

	fmt.Printf("\nGenerating geomap finished at %s !\n", time.Now().Format(time.RFC3339))

	data.export()
}

func getParsedList(parents []ParsedPair) []ParsedPair {
	var parsedList []ParsedPair

	res, err := http.Get(api + createParams(parents))
	if err != nil { panic(err) }
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != 200 { panic(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)) }

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil { panic(err) }

	selector := createSelector(parents)
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		parsedList = append(
			parsedList,
			ParsedPair{
				Code: value,
				Name: utils.EucKrToUtf8(s.Text()),
			},
		)
	})

	fmt.Print(".")
	return parsedList
}

// SI-DO = City(시) + Province(도)
// GU-GUN = District(구) + Country(군)
// DONG = Neighborhood(동)
const (
	sidoType = 0 + iota
	gugunType
	dongType
)

func createParams(args []ParsedPair) string {
	switch len(args) {
	case sidoType:
		return ""
	case gugunType:
		return fmt.Sprintf("?sido=%s", args[sidoType].Code)
	case dongType:
		return fmt.Sprintf("?sido=%s&gugun=%s", args[sidoType].Code, args[gugunType].Code)
	default:
		return ""
	}
}

func createSelector(args []ParsedPair) string {
	switch len(args) {
	case sidoType:
		return "#search_area option"
	case gugunType:
		return "#search_area2 option"
	case dongType:
		return "#search_area3 option"
	default:
		return "#search_area option"
	}
}