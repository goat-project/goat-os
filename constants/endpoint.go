package constants

// prefix for endpoint options
const cfgEndpointOptions = "endpoint-options."

const (
	// CfgEndpointType represents the service type for the client
	CfgEndpointType = cfgEndpointOptions + "type"

	// CfgEndpointName represents the service name for the client
	CfgEndpointName = cfgEndpointOptions + "name"

	// CfgEndpointRegion represents the geographic region in which the endpoint resides
	CfgEndpointRegion = cfgEndpointOptions + "region"

	// CfgEndpointAvailability represents the visibility of the endpoint to be returned
	CfgEndpointAvailability = cfgEndpointOptions + "availability"
)
