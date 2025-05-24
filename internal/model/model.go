package model

type Rate struct {
	InputSymbol  string  `json:"input_symbol"`
	OutputToken  string  `json:"output_token"`
	ProjectName  string  `json:"project_name"`
	PoolName     string  `json:"pool_name"`
	APY          float64 `json:"apy"`
	ProjectLink  string  `json:"project_link"`
}
