package app

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/nareix/curl"
	"github.com/olekukonko/tablewriter"
	toml "github.com/pelletier/go-toml"
)

var optSearch = &cli.OptionInfo{
	Name:      "search",
	UsageLine: "search [-type type][-string search]",
	Short:     "Search the Flogo Artifact Repository for activities and triggers",
	Long: `Search the Flogo Artifact Repository for activities and triggers.

Options:
    -type     the type you're looking for ("all"|"activity"|"trigger")
    -string   the search string you want to use
	
Example:
    flogo search -type activity -string "dynamo" would search the repository for activities connecting to DynamoDB
 `,
}

const (
	typeAll      = "all"
	typeActivity = "activity"
	typeTrigger  = "trigger"
	tomlURL      = "https://raw.githubusercontent.com/TIBCOSoftware/flogo/master/showcases/data/items.toml"
	tomlKey      = "items"
)

var (
	filterByType   = false
	filterByString = false
)

func init() {
	CommandRegistry.RegisterCommand(&cmdSearch{option: optSearch})
}

type cmdSearch struct {
	option       *cli.OptionInfo
	searchType   string
	searchString string
	refresh      bool
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdSearch) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdSearch) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.searchType), "type", "", "type")
	fs.StringVar(&(c.searchString), "string", "", "string")
}

// Exec implementation of cli.Command.Exec
func (c *cmdSearch) Exec(args []string) error {

	// More than 0 means more than the flags
	if len(args) > 0 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	// Break if the type is not known
	if len(c.searchType) > 0 {
		switch c.searchType {
		case typeAll,
			typeActivity,
			typeTrigger:
			filterByType = true
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown type - %s\n\n", c.searchType)
			cmdUsage(c)
		}
	}

	if len(c.searchString) > 0 {
		filterByString = true
	}

	if !filterByType && !filterByString {
		fmt.Fprint(os.Stderr, "Error: Neither type or string flags are specified\n\n")
		cmdUsage(c)
	}

	// Get the FAR content
	content, err := getFARContent()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		os.Exit(2)
	}

	// Find artifacts
	datamap, err := searchContent(content, filterByType, filterByString, c.searchType, c.searchString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		os.Exit(2)
	}

	// Print a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type", "Description", "URL", "Author"})

	for _, v := range datamap {
		table.Append(v)
	}

	table.Render()

	return nil
}

func searchContent(content string, filterByType bool, filterByString bool, searchType string, searchString string) ([][]string, error) {
	// Read the content
	config, err := toml.Load(content)
	if err != nil {
		return nil, err
	}

	// Get the correct key
	queryResult := config.Get(tomlKey)
	if queryResult == nil {
		return nil, fmt.Errorf("Unknown error occurred, no items found in Flogo Showcase")
	}

	// Prepare a result structure
	resultArray := queryResult.([]*toml.Tree)
	datamap := make([][]string, 0)
	for _, val := range resultArray {
		tempVal := val.ToMap()
		if filterByType && filterByString {
			if containsKey(tempVal, "type", searchType) && containsValue(tempVal, searchString) {
				datamap = append(datamap, []string{tempVal["name"].(string), tempVal["type"].(string), tempVal["description"].(string), tempVal["url"].(string), tempVal["author"].(string)})
			}
		} else if filterByType {
			if containsKey(tempVal, "type", searchType) {
				datamap = append(datamap, []string{tempVal["name"].(string), tempVal["type"].(string), tempVal["description"].(string), tempVal["url"].(string), tempVal["author"].(string)})
			}
		} else if filterByString {
			if containsValue(tempVal, searchString) {
				datamap = append(datamap, []string{tempVal["name"].(string), tempVal["type"].(string), tempVal["description"].(string), tempVal["url"].(string), tempVal["author"].(string)})
			}
		}
	}

	return datamap, nil
}

func containsKey(datamap map[string]interface{}, key string, value string) bool {
	if _, ok := datamap[key]; ok {
		if datamap[key] == value {
			return true
		}
	}
	return false
}

func containsValue(datamap map[string]interface{}, value string) bool {
	for key := range datamap {
		if strings.Contains(datamap[key].(string), value) {
			return true
		}
	}
	return false
}

func getFARContent() (string, error) {
	// Create new request
	req := curl.Get(tomlURL)

	// Set timeouts
	// DialTimeout is the TCP Connection Timeout
	// Timeout is the Download Timeout
	req.DialTimeout(time.Second * 10)
	req.Timeout(time.Second * 30)

	// Specify a progress monitor, otherwise it doesn't work
	req.Progress(func(p curl.ProgressStatus) {}, time.Second)

	// Execute the request and return the result
	res, err := req.Do()
	if err != nil {
		return "", err
	}

	if res.StatusCode == 200 {
		return res.Body, nil
	}

	return "", fmt.Errorf("Unknown error occurred with HTTP Status Code %v", res.StatusCode)
}
