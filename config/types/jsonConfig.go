package basicTypes

type JSONConfig interface {
	HasExternalFile() bool
	GetJSONPath() string
}
