package dockle

type ScanResult struct {
	Summary summary  `json:"summary"`
	Details []detail `json:"details"`
}

type summary struct {
	Fatal int `json:"fatal"`
	Warn  int `json:"warn"`
	Info  int `json:"info"`
	Skip  int `json:"skip"`
	Pass  int `json:"pass"`
}

type detail struct {
	Code   string   `json:"code"`
	Title  string   `json:"title"`
	Level  string   `json:"level"`
	Alerts []string `json:"alerts"`
}
