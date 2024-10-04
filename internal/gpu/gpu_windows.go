package gpu

import (
	"github.com/jaypipes/ghw"
)

func GetGPUInfo() ([]GPU, error) {
	graphics, err := ghw.GPU()

	if err != nil {
		return nil, err
	}

	cards := graphics.GraphicsCards

	parsedCards := make([]GPU, len(cards))
	for c := range cards {
		card := cards[c]
		parsedCards[c] = GPU{
			Id:          string(rune(card.Index)),
			Class:       "display",
			Description: card.DeviceInfo.Product.ID,
			Product:     card.DeviceInfo.Product.Name,
			Vendor:      card.DeviceInfo.Vendor.Name,
			PhysicalID:  card.DeviceInfo.Address,
			BusInfo:     card.DeviceInfo.ProgrammingInterface.Name,
		}
	}
	return parsedCards, nil
}
