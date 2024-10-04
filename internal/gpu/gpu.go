package gpu

import (
	"github.com/jaypipes/ghw"
	"github.com/kr/pretty"
)

type GPU struct {
	Id            int    `json:"id"`
	Class         string `json:"class"`
	Description   string `json:"description"`
	Product       string `json:"product"`
	Vendor        string `json:"vendor"`
	PhysicalID    string `json:"physid"`
	BusInfo       string `json:"businfo"`
	Configuration string `json:"configuration"`
}

func GetGPUInfo() ([]*GPU, error) {
	graphics, err := ghw.GPU()

	if err != nil {
		return nil, err
	}

	cards := graphics.GraphicsCards

	parsedCards := make([]*GPU, len(cards))
	for c := range cards {
		card := cards[c]

		pretty.Log(card)
		parsedCards[c] = &GPU{
			Id:            card.Index,
			Class:         "display",
			Description:   card.DeviceInfo.Product.ID,
			Product:       card.DeviceInfo.Product.Name,
			Vendor:        card.DeviceInfo.Vendor.Name,
			PhysicalID:    card.DeviceInfo.Address,
			BusInfo:       card.DeviceInfo.ProgrammingInterface.Name,
			Configuration: "unknown",
		}
	}
	return parsedCards, nil
}
