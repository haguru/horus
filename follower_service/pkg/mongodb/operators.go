package mongodb

// CurrentDateOp sets the value of a field to current date, either as a Date or a Timestamp.
type CurrentDateOp struct {
	CurrentDate interface{} `bson:"$currentDate"`
}

// IncOp increments the value of the field by the specified amount.
type IncOp struct {
	Inc interface{} `bson:"$inc"`
}

// MaxOp only updates the field if the specified value is greater than the existing field value.
type MaxOp struct {
	Max interface{} `bson:"$max"`
}

// MinOp only updates the field if the specified value is less than the existing field value.
type MinOp struct {
	Min interface{} `bson:"$min"`
}

// MulOp multiplies the value of the field by the specified amount.
type MulOp struct {
	Mul interface{} `bson:"$mul"`
}

// RenameOp Renames a field.
type RenameOp struct {
	Rename interface{} `bson:"$rename"`
}

// SetOp Sets the value of a field in a document.
type SetOp struct {
	Set interface{} `bson:"$set"`
}

// Sets the value of a field if an update results in an insert of a document.
// Has no effect on update operations that modify existing documents.
type SetOnInsertOp struct {
	SetOnInsert interface{} `bson:"$setOnInsert"`
}

// UnsetOp removes the specified field from a document.
type UnsetOp struct {
	Unset interface{} `bson:"$unset"`
}
