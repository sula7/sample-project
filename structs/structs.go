package structs

type AuthToken struct {
	AccessExpiresAt  int64
	RefreshExpiresAt int64
	AccessToken      string
	RefreshToken     string
	AccessUuid       string
	RefreshUuid      string
}

type AccessDetails struct {
	AccessUuid string
	UserUUID   string
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Drone struct {
	ID               uint64  `json:"id"`
	Name             *string `json:"name"`
	Description      *string `json:"description"`
	Author           string  `json:"author"`
	UserUUID         string  `json:"-"`
	Frame            *string `json:"frame"`
	Motors           *string `json:"motors"`
	ESC              *string `json:"esc"`
	Propellers       *string `json:"propellers"`
	FPVCamera        *string `json:"fpv_camera"`
	VTX              *string `json:"vtx"`
	VTXAntenna       *string `json:"vtx_antenna"`
	RX               *string `json:"rx"`
	FlightController *string `json:"flight_controller"`
	ExtraEquipment   *string `json:"extra_equipment"`
	CreatedAt        *string `json:"created_at"`
	UpdatedAt        *string `json:"updated_at"`
}
