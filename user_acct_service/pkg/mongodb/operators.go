package mongodb


type CurrentDateOp struct {
	CurrentDate interface{} `bson:"$currentDate"`

}

type IncOp struct {
	Inc         interface{} `bson:"$inc"`
}

type MaxOp struct {
	Max         interface{} `bson:"$max"`
}

type MinOp struct {
	Min         interface{} `bson:"$min"`
}

type MulOp struct {
	Mul         interface{} `bson:"$mul"`
}

type RenameOp struct {
	Rename      interface{} `bson:"$rename"`
}

type SetOp struct {
	Set         interface{} `bson:"$set"`
}

type SetOnInsertOp struct {
	SetOnInsert interface{} `bson:"$setOnInsert"`
}

type UnsetOp struct {
	Unset         interface{} `bson:"$unset"`
}
