package connectivity

type DestinationResponse struct {
	SystemMetadata           SystemMetadata    `json:"systemMetadata"`
	DestinationConfiguration map[string]string `json:"destinationConfiguration"`
	PropertiesMetadata       []interface{}     `json:"propertiesMetadata"`
}

type SystemMetadata struct {
	CreationTime     string `json:"creation_time"`
	Etag             string `json:"etag"`
	UserAgent        string `json:"userAgent"`
	ModificationTime string `json:"modification_time"`
}

type DestinationName struct {
	Name string `json:"Name"`
}

type DestinationNamesList struct {
	Names []DestinationName `json:"Names"`
}

type DestinationsList struct {
	Destination []DestinationResponse `json:"destination"`
}
