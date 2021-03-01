package dummy

import "github.com/go-openapi/strfmt"

// TDummy defines the type of a dummy structure
// swagger:model
// unique: true
type TDummy struct {
	// required: true
	ID strfmt.UUID `json:"id" validate:"required"` // Unique identifier

	// max length: 100
	Description string `json:"description"`

	// swagger:ignore
	Action string `json:"action,omitempty"`
}

// TDummies defines the type of a dummySS list
type TDummies []*TDummy

/*
// dummyList defines a list with some dummies. It should be replaced by a database approach.
var dummyList = TDummies{
	&TDummy{
		//ID:          uuid.Must(uuid.Parse("ff240753-74b0-4a82-8433-275aef1f3b55")),
		ID:          "ff240753-74b0-4a82-8433-275aef1f3b55",
		Description: "Desc 1",
	},
	&TDummy{
		//ID:          uuid.Must(uuid.Parse("8a5496ab-da47-4015-bb76-4edc4768e199")),
		ID:          "8a5496ab-da47-4015-bb76-4edc4768e199",
		Description: "Desc 2",
	},
}
*/
