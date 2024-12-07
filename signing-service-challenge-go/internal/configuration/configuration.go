package configuration

// Configuration will hold our internal configuration settings
type Configuration struct {
	ListenAddress string `json:"listen_address"`
}

// LoadConfiguration in real live we would load the env file here, some other way of getting the env variables
func LoadConfiguration() (*Configuration, error) {
	return &Configuration{
		ListenAddress: ":8080",
	}, nil
}
