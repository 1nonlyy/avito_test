package models

type PVZWithReceptions struct {
	PVZ        PVZ                  `json:"pvz"`
	Receptions []ReceptionWithItems `json:"receptions"`
}

type ReceptionWithItems struct {
	Reception Reception `json:"reception"`
	Products  []Product `json:"products"`
}
