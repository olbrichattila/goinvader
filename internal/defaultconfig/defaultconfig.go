// Package defaultconfig only contains default variables
package defaultconfig

type AlibabaServiceType int

// Default variables
const (
	ScreenW = 840
	ScreenH = 480
	// ApiUrl  = "http://localhost:3000/"
	ApiUrl = "https://www.alirobo.fun/api/"
)

// Events
const (
	_ AlibabaServiceType = iota
	BossRom
	Ecs
	FunctionCompute
	ServerlessComputing
	ObjectStorageService
	BlockStorage
	CloudBackup
	Cdn
	ApsaraDB
)

var ServiceDescriptionMap = map[AlibabaServiceType]string{
	Ecs:                  "ECS Elastic container service",
	FunctionCompute:      "Function compute",
	ServerlessComputing:  "Serverless Compute",
	ObjectStorageService: "OSS Object Storage Service",
	BlockStorage:         "Block Storage",
	CloudBackup:          "Cloud Backup",
	Cdn:                  "CDN",
	ApsaraDB:             "ApsaraDB Database storage",
}
