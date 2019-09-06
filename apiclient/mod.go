package apiclient

import (
	"fmt"
	"log"
	"strings"
)

func (client *Client) InstallMod(parent, org, mod, version string) (*TurbotMetadata, error) {
	query := installModMutation()
	responseData := &InstallModResponse{}

	commandPayload := map[string]string{
		"parentAka": parent,
		"org":       org,
		"mod":       mod,
		"version":   version,
	}

	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error installing mod: %s", err.Error())
	}

	return &responseData.Mod.Turbot, nil
}

func (client *Client) ReadMod(id string) (*Mod, error) {
	query := readModQuery(id)
	responseData := &ReadModResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading policy: %s", err.Error())
	}

	log.Println("ReadMod", responseData)
	// convert uri into org and mod
	org, mod := ParseModUri(responseData.Mod.Uri)
	var result = &Mod{
		Org:     org,
		Mod:     mod,
		Version: responseData.Mod.Version,
		Parent:  responseData.Mod.Parent,
	}
	return result, nil
}

func ParseModUri(uri string) (org, mod string) {
	if uri == "" {
		org = ""
		mod = ""
		return
	}
	// uri will be of form "tmod:@<org>/<mod>"
	segments := strings.Split(strings.TrimPrefix(uri, "tmod:@"), "/")
	org = segments[0]
	mod = segments[1]
	return
}

func (client *Client) UninstallMod(modId string) error {
	query := uninstallModMutation()
	responseData := &UninstallModResponse{}

	commandPayload := map[string]string{
		"modResourceId": modId,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
		},
	}

	log.Println("UninstallMod")
	log.Println("modId", modId)
	log.Println("Query:", query)

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error uninstalling mod: %s", err.Error())
	}
	log.Println("DoRequest returned:", responseData)
	if !responseData.ModUninstall.Success {
		return fmt.Errorf("modUninstall mutation ran with no errors but failed to uninstall the mod")
	}
	return nil
}
