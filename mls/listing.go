package mls

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var fileModePrefix = "_tmp/"

type Listing struct {
	MLSID         string
	State         string
	Kind          string
	Address       string
	ListPrice     int
	HOAFees       int
	TotalRooms    int
	Bedrooms      int
	FullBaths     int
	HalfBaths     int
	MasterBath    bool
	UnitLevel     int
	LivingArea    int
	PricePerSqF   int
	ParkingSpaces int
	AssessedValue int
	AssessedYear  int
}

func GetMLSIDs() ([]string, error) {

	ids := []string{}
	var data *goquery.Document
	mlsSet := make(map[string]bool)

	if Config.FileMode {
		f, err := os.Open(fileModePrefix + "listing.html")
		if err != nil {
			return ids, err
		}

		data, err = goquery.NewDocumentFromReader(f)
		if err != nil {
			return ids, err
		}

	} else {
		// get from web
		resp, err := webGet(Config.BaseUrl + "/index.aspx")
		if err != nil {
			return ids, err
		}
		defer resp.Body.Close()

		data, err = goquery.NewDocumentFromResponse(resp)
		if err != nil {
			return ids, err
		}

		if Config.DumpPages {
			body, _ := data.Html()
			if err := ioutil.WriteFile(fileModePrefix+"listing.html", []byte(body), 0644); err != nil {
				return ids, err
			}
		}

	}
	data.Find(".VOWResultsRow a").Each(
		func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			u, _ := url.Parse(href)

			if m := u.Query()["mls"]; m != nil {
				mlsSet[m[0]] = true
			}
		},
	)

	for m, _ := range mlsSet {
		ids = append(ids, m)
	}

	return ids, nil
}

func GetListing(id string) (Listing, error) {
	l := Listing{MLSID: id}
	var data *goquery.Document
	if Config.FileMode {
		f, err := os.Open(fileModePrefix + id + ".html")
		if err != nil {
			return l, err
		}

		data, err = goquery.NewDocumentFromReader(f)
		if err != nil {
			return l, err
		}

	} else {

		url := Config.BaseUrl + "/reports/public.aspx?ShowRatePlug=False&mls=" + id
		resp, err := webGet(url)
		if err != nil {
			return l, err
		}
		defer resp.Body.Close()

		data, err = goquery.NewDocumentFromResponse(resp)
		if err != nil {
			return l, err
		}
	}

	if Config.DumpPages {
		body, err := data.Html()
		if err != nil {
			return l, err
		}
		if err := ioutil.WriteFile(fileModePrefix+id+".html", []byte(body), 0644); err != nil {
			return l, err
		}
	}

	// Process the data

	// State, Type
	data.Find("body table tbody tr td table tbody tr td table tbody tr td table tbody tr td.big b").Each(
		func(i int, s *goquery.Selection) {
			a := strings.Split(s.Text(), "\n")
			re := regexp.MustCompile(`\w+$`)
			l.State = re.FindString(a[3])
			l.Kind = re.FindString(a[4])
		},
	)

	re := regexp.MustCompile(`((\w+\s?)+|\w+):\s+([$|\w|,|\w\s]+)\n`)
	reDigit := regexp.MustCompile(`[^0-9]`)
	reBath := regexp.MustCompile(`(\d+)f\s(\d+)h`)

	data.Find("body table tbody tr td table tbody tr td table tbody tr td").Each(
		func(i int, s *goquery.Selection) {

			matches := re.FindStringSubmatch(s.Text())
			debug := "71734795" == id
			if debug {
				log.Printf("%d: %q", i, s.Text())
			}
			if len(matches) > 1 {

				if debug {
					log.Printf("%s => %s\n", matches[1], matches[3])
				}
				switch matches[1] {
				case "Bedrooms":
					l.Bedrooms, _ = strconv.Atoi(matches[3])
				case "Total Rooms":
					l.TotalRooms, _ = strconv.Atoi(matches[3])
				case "Master Baths":
					if matches[2] == "Yes" {
						l.MasterBath = true
					}
				case "Unit Level":
					l.UnitLevel, _ = strconv.Atoi(matches[3])
				case "Living Area":
					l.LivingArea, _ = strconv.Atoi(matches[3])
				case "List Price":
					l.ListPrice, _ = strconv.Atoi(reDigit.ReplaceAllString(matches[3], ""))
				case "Bathrooms":
					a := reBath.FindStringSubmatch(matches[3])
					l.FullBaths, _ = strconv.Atoi(a[1])
					l.HalfBaths, _ = strconv.Atoi(a[2])
				case "Assessed":
					l.AssessedValue, _ = strconv.Atoi(reDigit.ReplaceAllString(matches[3], ""))
				case "Tax Year":
					l.AssessedYear, _ = strconv.Atoi(matches[3])
				}
			} else {
				//log.Printf("NOMATCH\n")
			}
		},
	)
	return l, nil
}
