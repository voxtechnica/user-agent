package main

import (
	"encoding/json"
	"fmt"
	"github.com/voxtechnica/user-agent"
	"log"
	"os"
	"sort"
	"strings"
)

type userAgentCount struct {
	UserAgent     string         `json:"userAgent"`
	Count         int            `json:"count"`
	VersionCounts map[string]int `json:"versionCounts"`
	StringCounts  map[string]int `json:"stringCounts"`
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Printf("error closing file %s: %v", f.Name(), err)
	}
}

func main() {
	// load user agent test data
	var userAgentCounts = make([]userAgentCount, 212)
	sampleData, err := os.Open("sample_data/user_agents.json")
	if err != nil {
		log.Fatalln("error reading test data:", err)
	}
	defer closeFile(sampleData)
	err = json.NewDecoder(sampleData).Decode(&userAgentCounts)
	if err != nil {
		log.Fatalln("error decoding JSON test data:", err)
	}

	// open results files
	output, err := os.Create("sample_data/user_agents.txt")
	if err != nil {
		log.Fatalln("error opening file user_agents.txt:", err)
	}
	defer closeFile(output)

	// parse user agent strings
	uaCount := 0
	uaStringCount := 0
	clientTypeCounts := map[string]int{}
	clientNameCounts := map[string]int{}
	deviceTypeCounts := map[string]int{}
	osNameCounts := map[string]int{}
	urlCounts := map[string]int{}
	for _, uac := range userAgentCounts {
		uaStringCount += len(uac.StringCounts)
		for s, c := range uac.StringCounts {
			uaCount += c
			userAgent := user_agent.Parse(s)
			ctCount := clientTypeCounts[userAgent.ClientType]
			clientTypeCounts[userAgent.ClientType] = ctCount + c
			cnCount := clientNameCounts[userAgent.ClientName]
			clientNameCounts[userAgent.ClientName] = cnCount + c
			dtCount := deviceTypeCounts[userAgent.DeviceType]
			deviceTypeCounts[userAgent.DeviceType] = dtCount + c
			osCount := osNameCounts[userAgent.OSName]
			osNameCounts[userAgent.OSName] = osCount + c
			if userAgent.URL != "" {
				urlCount := urlCounts[userAgent.URL]
				urlCounts[userAgent.URL] = urlCount + c
			}
			cleaned := strings.Join(userAgent.Fields, " ")
			parsed := userAgent.String()
			uas := fmt.Sprintf("%s\n%s\t(%d occurrences)\n%s\n\n", userAgent.Header, cleaned, c, parsed)
			_, err = output.WriteString(uas)
			if err != nil {
				log.Printf("error writing to file %s: %v", output.Name(), err)
			}
		}
	}
	err = output.Sync()
	if err != nil {
		log.Printf("error syncing file %s: %v", output.Name(), err)
	}

	fmt.Println("User-Agent Count:", len(userAgentCounts))
	fmt.Println("User-Agent String Count:", uaStringCount)
	fmt.Println("User-Agent View Count:", uaCount)

	printCounts(clientTypeCounts, "Client Type")
	printCounts(clientNameCounts, "Client Name")
	printCounts(deviceTypeCounts, "Device Type")
	printCounts(osNameCounts, "OS Name")
	printCounts(urlCounts, "URL")
}

func printCounts(counts map[string]int, title string) {
	fmt.Println("\n Views\tPercent\t" + title)
	var keys []string
	var total int
	for k, v := range counts {
		keys = append(keys, k)
		total += v
	}
	sort.Strings(keys)
	for _, k := range keys {
		p := float32(counts[k]) / float32(total) * 100
		fmt.Printf("%6d\t%6.2f\t%s\n", counts[k], p, k)
	}
	fmt.Printf("%6d\t100.00\tTotal\n", total)
}
