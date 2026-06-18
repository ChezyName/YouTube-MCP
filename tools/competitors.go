package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/ChezyName/YouTube-MCP/config"
)

var competitors *[]Competitor

func IsCompetitorsEnabled() bool {
	cfg := config.GetConfig()
	if cfg == nil {
		return false
	}

	return cfg.Competitors
}

func loadCompetitors() {
	appData, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR Trying to look for appdata path:%s", err.Error())
		os.Exit(1)
	}

	configDir := filepath.Join(appData, "/YouTube-MCP")
	competitorsFile := filepath.Join(configDir, config.CompetitorsFile)

	//make the folder if does not exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}

	//does not exist
	if _, err := os.Stat(competitorsFile); os.IsNotExist(err) {
		defaultCfg := []Competitor{}

		bytes, _ := json.MarshalIndent(defaultCfg, "", "  ")

		err := os.WriteFile(competitorsFile, bytes, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not create competitors file: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "A default competitors file has been created at:\n%s\n", competitorsFile)
	}

	file, err := os.ReadFile(competitorsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not read competitors file: %v\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(file, &competitors); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Invalid JSON in competitors file: %v\n", err)
		os.Exit(1)
	}
}

func saveCompetitors() error {
	if competitors == nil {
		return fmt.Errorf("Cannot save an empty or nil competitors object")
	}
	bytes, _ := json.MarshalIndent(competitors, "", "  ")

	appData, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR SAVING COMPETITORS: Trying to look for appdata path:%s", err.Error())
		os.Exit(1)
	}

	configDir := filepath.Join(appData, "/YouTube-MCP")
	competitorsFile := filepath.Join(configDir, config.CompetitorsFile)

	err = os.WriteFile(competitorsFile, bytes, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR SAVING COMPETITORS: Could not create competitors file: %v\n", err)
		os.Exit(1)
	}

	return nil
}

func GetCompetitors() ([]Competitor, error) {
	if competitors == nil {
		loadCompetitors()
	}

	if competitors == nil {
		return nil, fmt.Errorf("Failed to load competitors file properly")
	}

	return *competitors, nil
}

func AddCompetitor(name string, tags []string) error {
	if name == "" {
		return fmt.Errorf("AddCompetitor needs a channel name / handle.")
	}

	cmp, err := GetCompetitors()
	if err != nil {
		return err
	}

	//check for dupes:
	var foundDupe = -1
	for i, v := range cmp {
		if v.Name == name {
			foundDupe = i
			break
		}
	}

	if foundDupe == -1 {
		//add
		cmp = append(cmp, Competitor{
			Name: name,
			Tags: tags,
		})

	} else {
		fmt.Fprintf(os.Stderr, "Competitor %s already in list of competitors, going to update tags", name)
		cmp[foundDupe].Tags = slices.Concat(cmp[foundDupe].Tags, tags)
	}

	//save to config file
	competitors = &cmp
	err = saveCompetitors()
	return err
}

func RemoveCompetitor(name string) error {
	if name == "" {
		return fmt.Errorf("RemoveCompetitor needs a channel name / handle.")
	}

	cmp, err := GetCompetitors()
	if err != nil {
		return err
	}

	//check for dupes:
	var foundDupe = -1
	for i, v := range cmp {
		if v.Name == name {
			foundDupe = i
			break
		}
	}

	if foundDupe == -1 {
		return fmt.Errorf("Could not find competitor")
	} else {
		//add
		cmp = append(cmp[:foundDupe], cmp[foundDupe+1:]...)
	}

	//save to config file
	competitors = &cmp
	err = saveCompetitors()
	return err
}
