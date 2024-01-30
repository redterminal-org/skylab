package config

type Config struct {
	GeminiFetchAddress	string
	CertFilename    	string
	AutoFetch			bool
	Timezone			string
	// Add other configuration fields as needed
}

func ExportConfig() *Config {
    // Configuration
    return &Config{
        GeminiFetchAddress: "gemini://satch.xyz:1958/fetch",
        CertFilename:       "../cert.pem",
		AutoFetch:			false,
		Timezone:			"America/New_York",
        // Add other configuration fields as needed
    }
}
