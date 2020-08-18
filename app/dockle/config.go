package dockle

type Config struct {
	ScanImageName                string
	IsLocalImage                 bool
	DockleAuthUrl                string
	DockleUsername               string
	DocklePassword               string
	AwsAccessKeyId               string
	AwsSecretAccessKey           string
	AwsDefaultRegion             string
	GoogleApplicationCredentials string
	GithubToken                  string
}
