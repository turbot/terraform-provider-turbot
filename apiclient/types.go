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
	Turbot             TurbotPolicyMetadata
}

type TurbotResourceMetadata struct {
	Id       string
	ParentId string
	Akas     []string
}

type TurbotPolicyMetadata struct {
	Id       string
	ParentId string
	Akas     []string
}

type TurbotGrantMetadata struct {
	Id         string
	ProfileId  string
	ResourceId string
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
	Turbot     TurbotPolicyMetadata
}

// InstallModResponse: must be consistent with with fields defined in installModMutation
type InstallModResponse struct {
	Mod InstallModData
}

type InstallModData struct {
	Build  string
	Turbot TurbotResourceMetadata
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
		Turbot TurbotResourceMetadata
	}
}

type CreateGrantResponse struct {
	Grants struct {
		Items []struct {
			Turbot TurbotGrantMetadata
		}
	}
}

type UpdateResourceResponse struct {
	Resource struct {
		Turbot TurbotResourceMetadata
	}
}

// note: the Resource property is just an interface{} - this is mapped manually into a Resource object,
// rather than unmarshalled. This is to allow for dynamic data types, while always having the Turbot property
type ReadResourceResponse struct {
	Resource interface{}
}

type ReadFullResourceResponse struct {
	Resource FullResource
}

type ReadResourceListResponse struct {
	ResourceList struct {
		Items []Resource
	}
}

type Resource struct {
	Turbot TurbotResourceMetadata
	Data   map[string]interface{}
}

type FullResource struct {
	Object interface{}
	Turbot TurbotResourceMetadata
}

type ReadFolderResponse struct {
	Resource Folder
}

type ReadSmartFolderResponse struct {
	Resource SmartFolder
}

type Profile struct {
	Turbot          TurbotResourceMetadata
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

type ReadProfileResponse struct {
	Resource Profile
}

type FindFolderResponse struct {
	Folders struct {
		Items []Folder
	}
}

type FindSmartFolderResponse struct {
	SmartFolders struct {
		Items []SmartFolders
	}
}

type Folder struct {
	Turbot      TurbotResourceMetadata
	Title       string
	Description string
	Parent      string
}

type SmartFolder struct {
	Turbot      TurbotMetadata
	Title       string
	Description string
	Rules        map[string]interface{}
	Parent      string
}

type ReadLocalDirectoryResponse struct {
	Resource LocalDirectory
}

type LocalDirectory struct {
	Turbot            TurbotResourceMetadata
	Title             string
	Description       string
	Parent            string
	Status            string
	DirectoryType     string
	ProfileIdTemplate string
}

type ReadSamlDirectoryResponse struct {
	Resource SamlDirectory
}

type SamlDirectory struct {
	Turbot            TurbotResourceMetadata
	Title             string
	Description       string
	Parent            string
	Status            string
	DirectoryType     string
	ProfileIdTemplate string
	EntryPoint        string
	Certificate       string
}
type ReadLocalDirectoryUserResponse struct {
	Resource LocalDirectoryUser
}

type LocalDirectoryUser struct {
	Turbot      TurbotResourceMetadata
	Parent      string
	Title       string
	Email       string
	Status      string
	DisplayName string
	GivenName   string
	MiddleName  string
	FamilyName  string
	Picture     string
}

type ReadGoogleDirectoryResponse struct {
	Resource GoogleDirectory
}

type GoogleDirectory struct {
	Turbot            TurbotResourceMetadata
	Parent            string
	Title             string
	ProfileIdTemplate string
	Description       string
	Status            string
	DirectoryType     string
	ClientID          string
	ClientSecret      string
	PoolId            string
	GroupIdTemplate   string
	LoginNameTemplate string
	HostedName        string
}

type ReadGrantResponse struct {
	Grant Grant
}

type Grant struct {
	Turbot            TurbotGrantMetadata
	PermissionTypeId  string
	PermissionLevelId string
}
