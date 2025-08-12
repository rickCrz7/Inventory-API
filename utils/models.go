package utils

type Owner struct {
	ID        string  `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	CampusID  *string `json:"campusID"`
	Email     string  `json:"email"`
}

type Type struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type TypeProperty struct {
	ID       string `json:"id"`
	TypeID   string `json:"typeID"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
	Required bool   `json:"required"`
}

type Device struct {
	ID           string `json:"id"`
	SerialNumber string `json:"serialNumber"`
	Name         string `json:"name"`
	TypeID       string `json:"typeID"`
	OwnerID      string `json:"ownerID"`
	PurchaseDate string `json:"purchaseDate"`
	Status       string `json:"status"`
}

type DeviceProperty struct {
	ID              string `json:"id"`
	DeviceID        string `json:"deviceID"`
	TypePropertyID  string `json:"typePropertyID"`
	Value           string `json:"value"`
}

type DevicePhoto struct {
	ID        string `json:"id"`
	DeviceID  string `json:"deviceID"`
	Photo     string `json:"photo"`
	CreatedAt string `json:"createdAt"`
}

type DeviceLog struct {
	ID        string `json:"id"`
	DeviceID  string `json:"deviceID"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}
