package gpu

type GPU struct {
	Id          string `json:"id"`
	Class       string `json:"class"`
	Description string `json:"description"`
	Product     string `json:"product"`
	Vendor      string `json:"vendor"`
	PhysicalID  string `json:"physid"`
	BusInfo     string `json:"businfo"`
}
