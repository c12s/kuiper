package model

type Label struct {
	Key   string
	Value string
}

type Config struct {
	Key    string
	Value  string
	Labels []Label
}

type Group struct {
	Configs []Config
	Version string
}
