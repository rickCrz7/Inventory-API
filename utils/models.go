package utils

type Owner struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	CampusID  *string `json:"campus_id"`
	Email     string  `json:"email"`
}

type Type struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type TypeProperty struct {
	ID       string `json:"id"`
	TypeID   string `json:"type_id"`
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	Required bool   `json:"required"`
}

type Device struct {
	ID           string `json:"id"`
	SerialNumber string `json:"serial_number"`
	Name         string `json:"name"`
	TypeID       string `json:"type_id"`
	OwnerID      string `json:"owner_id"`
	PurchaseDate string `json:"purchase_date"`
	Status       string `json:"status"`
}

type DeviceProperty struct {
	ID              string `json:"id"`
	DeviceID        string `json:"device_id"`
	TypePropertyID  string `json:"type_property_id"`
	Value           string `json:"value"`
}

type DevicePhoto struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id"`
	Photo     string `json:"photo"`
	CreatedAt string `json:"created_at"`
}

type DeviceLog struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
