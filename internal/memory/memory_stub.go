package memory

const (
	GIGABYTES = 1073741824
)

type MemoryInfo struct {
	Id     string `json:"id"`
	Class  string `json:"class"`
	PhysId string `json:"physid"`
	Units  string `json:"units"`
	Size   int64  `json:"size"`
	Vendor string `json:"vendor"`
	// Currently cannot get module speed
	Speed      int64 `json:"speed"`
	FormFactor string
}

//func Get_memory_info() ([]*MemoryInfo, error) {
//	return nil, errors.New("Method not found for OS: " + runtime.GOOS)
//}
