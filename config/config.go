package config

type Config struct {
	GeminiFetchAddress string
	CertFilename       string
	// Add other configuration fields as needed
}

func ExportConfig() *Config {
    // Configuration
    return &Config{
        GeminiFetchAddress: "gemini://satch.xyz:1958/fetch",
        CertFilename:       "../cert.pem",
        // Add other configuration fields as needed
    }
}
