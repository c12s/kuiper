package store

type ConfigGroupDAO struct {
	OrgId      string
	Name       string
	Version    string
	ParamsSets []struct {
		Name     string
		ParamSet map[string]string
	}
}