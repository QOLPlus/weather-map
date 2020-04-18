package geomap

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"net/http"
)

const api = "https://www.weather.go.kr/weather/lifenindustry/sevice_rss.jsp"

type ParsedPair struct {
	Code string
	Name string
}
func (p ParsedPair) String() string {
	return fmt.Sprintf("%s : %s", p.Name, p.Code)
}

// SI-DO = City(시) + Province(도)
// GU-GUN = District(구) + Country(군)
// DONG = Neighborhood(동)
const (
	sidoType = 0 + iota
	gugunType
	dongType
)

func Process(cmd *cobra.Command, args []string) {
	for _, sido := range getParsedList([]ParsedPair{}) {
		fmt.Println(sido.String())
		for _, gugun := range getParsedList([]ParsedPair{sido}) {
			fmt.Println("  ", gugun.String())
			for _, dong := range getParsedList([]ParsedPair{sido, gugun}) {
				fmt.Println("    ", dong.String())
			}
		}
	}
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
				Name: eucKrToUtf8(s.Text()),
			},
		)
	})

	return parsedList
}

func eucKrToUtf8(s string) string {
	var buffers bytes.Buffer
	tr := transform.NewWriter(&buffers, korean.EUCKR.NewDecoder())
	defer func() { _ = tr.Close() }()
	_, _ = tr.Write([]byte(s))
	return buffers.String()
}

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