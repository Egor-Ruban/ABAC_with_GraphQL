package model

import "time"

type Environment struct {
	Date *time.Time `json:"Date"`
}

type Object struct {
	ID      *string `json:"ID" db:"object_id"`
	Name    *string `json:"Name" db:"name"`
	OwnerID *string `json:"Owner_ID" db:"owner_id"`
}

type Subject struct {
	ID   *string `json:"ID" db:"subject_id"`
	Name *string `json:"Name" db:"name"`
	Age  *int    `json:"Age" db:"age"`
	Role *Role   `json:"Role" db:"role"`
}
