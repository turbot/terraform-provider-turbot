package apiclient

// Validation response - returned by API validation call
type ValidationResponse struct {
	Schema struct {
		QueryType struct {
			Name string
		}
	}
}

// is the validation response successful?
func (response *ValidationResponse) isValid() bool {
	return response.Schema.QueryType.Name == "Query"
}

// ApiResponse: used to unmarshall API error responses
type ApiResponse struct {
	Errors []Error
}
type Error struct {
	Message string
}

// PolicySettingResponse: must be consistent with fields defined in readPolicySettingQuery
type PolicySettingResponse struct {
	PolicySetting PolicySetting
}

type FindPolicySettingResponse struct {
	PolicySettings struct {
		Items []PolicySetting
	}
}

type PolicySetting struct {
	Value              interface{}
	ValueSource        string
	Default            bool
	Precedence         string
	Template           string
	TemplateInput      string
	Input              string
	Note               string
	ValidFromTimestamp string
	ValidToTimestamp   string
	Turbot             TurbotMetadata
}

type TurbotMetadata struct {
	Id       string
	ParentId string
	Akas     []string
}

// PolicyValueResponse: must be consistent with fields defined in readPolicyValueQuery
type PolicyValueResponse struct {
	PolicyValue PolicyValue
}

type PolicyValue struct {
	Value      interface{}
	Precedence string
	State      string
	Reason     string
	Details    string
	Setting    PolicySetting
	Turbot     TurbotMetadata
}

// InstallModResponse: must be consistent with with fields defined in installModMutation
type InstallModResponse struct {
	Mod InstallModData
}

type InstallModData struct {
	Build  string
	Turbot TurbotMetadata
}

type ReadModResponse struct {
	Mod Mod
}

type ModRegistryVersion struct {
	Status  string
	Version string
}

type ModVersionResponse struct {
	Versions struct {
		Items []ModRegistryVersion
	}
}

type UninstallModResponse struct {
	ModUninstall struct {
		Success bool
	}
}

type Mod struct {
	Org     string
	Mod     string
	Version string
	Parent  string
	Uri     string
}

type CreateResourceResponse struct {
	Resource struct {
		Turbot TurbotMetadata
	}
}

type UpdateResourceResponse struct {
	Resource struct {
		Turbot TurbotMetadata
	}
}

// note: the Resource property is just an interface{} - this is mapped manually into a Resource object,
// rather than unmarshalled. This is to allow for dynamic data types, while always having the Turbot property
type ReadResourceResponse struct {
	Resource interface{}
}

type Resource struct {
	Turbot TurbotMetadata
	Data   map[string]interface{}
}

type ReadFolderResponse struct {
	Resource Folder
}

type Profile struct {
	Turbot          TurbotMetadata
	Title           string
	Parent          string
	Status          string
	Email           string
	GivenName       string
	DisplayName     string
	FamilyName      string
	DirectoryPoolId string
	ProfileId       string
}

type ProfilePayload struct {
	Title           string
	Parent          string
	Status          string
	DisplayName     string
	Email           string
	GivenName       string
	FamilyName      string
	DirectoryPoolId string
	ProfileId       string
}

type ReadProfileResponse struct {
	Resource Profile
}

type FindFolderResponse struct {
	Folders struct {
		Items []Folder
	}
}

type Folder struct {
	Turbot      TurbotMetadata
	Title       string
	Description string
	Parent      string
}
