package main

type tomlConfig struct {
	Title   string
	Owner   ownerInfo
	Token   token
	Servers map[string]server
	Files   map[string]files
	Hotkey  hotkey
	User    user
	Logs    logs
	Winscp  winscp
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
}

type token struct {
	TokenWait int
}

type user struct {
	Name string
}

type logs struct {
	Location   string
	Name       string
	Flag       int
	Permission int
}

type winscp struct {
	Location       string
	Proxy          string
	OutputLocation string
	Command        string
	TransFile      string
	Log            string
	LogLoc         string
}

type hotkey struct {
	Command        string
	Location       string
	OutputLocation string
}

type server struct {
	URL  string
	IP   string
	Port string
}

type files struct {
	ServerLocation string
	ServerName     string
	Get            bool
	OutputLocation string
	OutputName     string
	Prefix         string
}

/* type clients struct {
	Data  [][]interface{}
	Hosts []string
} */
