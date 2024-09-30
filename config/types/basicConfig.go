package basicTypes

type BasicConfig struct {
	RewriteWithJSON bool `json:"-"` // Enables JSON configuration. If it exists, this configuration will be shadowed by JSON one
}

func (b BasicConfig) HasExternalFile() bool {
	return b.RewriteWithJSON
}
