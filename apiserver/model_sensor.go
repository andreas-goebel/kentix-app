/*
 * Kentix app API
 *
 * API to access and configure the Kentix app
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// Sensor - Each sensor represents one asset in Eliona.
type Sensor struct {

	// The project ID this asset is assigned to
	ProjectID string `json:"projectID,omitempty"`

	Configuration Configuration `json:"configuration,omitempty"`

	// Eliona asset ID
	AssetID *int32 `json:"assetID,omitempty"`

	// Serial number reported by the Kentix device
	SerialNumber string `json:"serialNumber,omitempty"`
}

// AssertSensorRequired checks if the required fields are not zero-ed
func AssertSensorRequired(obj Sensor) error {
	if err := AssertConfigurationRequired(obj.Configuration); err != nil {
		return err
	}
	return nil
}

// AssertRecurseSensorRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Sensor (e.g. [][]Sensor), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseSensorRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aSensor, ok := obj.(Sensor)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertSensorRequired(aSensor)
	})
}
