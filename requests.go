package goSmartSheet

//ContainerDestination represents an destination container target for a copy operation
//https://smartsheet-platform.github.io/api-docs/#containerdestination-object
type ContainerDestination struct {
	Type          DestinationType `json:"destinationType"`
	DestinationID int64           `json:"destinationId,omitempty"`
	NewName       string          `json:"newName,omitempty"`
}

//DestinationType represents the possible destination types for a ContainerDestination
type DestinationType string

const (
	DestinationTypeHome     DestinationType = "home"
	DestinatonTypeWorkspace                 = "workspace"
	DestinatonTypeFolder                    = "folder"
)
