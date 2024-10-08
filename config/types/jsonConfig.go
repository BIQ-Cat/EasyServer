package basictypes

type JSONConfig interface {
	HasExternalFile() bool
	GetJSONPath() string
}
