package mongodb

import "fmt"

const (
	POINT_TYPE_POLYGON       = "Polygon"
	POINT_TYPE_MULTI_POLYGON = "MultiPolygon"
	POINT_TYPE_MULTI_POINT   = "Point"

	OP_TYPE_GEO_INTERSECTS = "geoIntersects"
	OP_TYPE_GEO_WITHIN     = "geoWithin"
	OP_TYPE_NEAR           = "near"
	OP_TYPE_NEAR_SPHERE    = "nearSphere"
)

type SpatialQueryCommand struct {
	Location interface{} `bson:"location"`
}

type GeoIntersectsOp struct {
	GeoIntersects interface{} `bson:"$geoIntersects"`
}

type Point struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type GeometryOP struct {
	Geometry    Point `bson:"$geometry"`
	MaxDistance int   `bson:"$maxDistance,omitempty"`
	MinDistance int   `bson:"$minDistance,omitempty"`
}

type GeoWithinOP struct {
	GeoWithin interface{} `bson:"$geoWithin"`
}

type NearOp struct {
	Near interface{} `bson:"$near"`
}

type NearSphereOP struct {
	NearSphere interface{} `bson:"$nearSphere"`
}

type MaxDistanceOp struct {
	MaxDistance int `bson:"$maxDistance"`
}

type MinDistanceOp struct {
	MinDistance int `bson:"$minDistance"`
}

func NewSpatialQueryCommand(opType string, pointType string, coordinates []float64, maxDistance int, minDistance int) (interface{}, error) {
	cmd := SpatialQueryCommand{}
	geometryOp := GeometryOP{}
	geometryOp.Geometry.Coordinates = coordinates
	// verify point type is allowed
	switch pointType {

	case POINT_TYPE_POLYGON:
		geometryOp.Geometry.Type = POINT_TYPE_POLYGON

	case POINT_TYPE_MULTI_POLYGON:
		geometryOp.Geometry.Type = POINT_TYPE_MULTI_POLYGON

	case POINT_TYPE_MULTI_POINT:
		if opType == OP_TYPE_GEO_WITHIN {
			return nil, fmt.Errorf("point type %v not allowed", pointType)
		}
		geometryOp.Geometry.Type = POINT_TYPE_MULTI_POINT
	default:
		return nil, fmt.Errorf("point type %v not allowed", pointType)
	}

	switch opType {
	case OP_TYPE_GEO_INTERSECTS:
		cmd.Location = GeoIntersectsOp{
			GeoIntersects: geometryOp,
		}
	case OP_TYPE_GEO_WITHIN:
		cmd.Location = GeoWithinOP{
			GeoWithin: geometryOp,
		}
	case OP_TYPE_NEAR:
		geometryOp.MaxDistance = maxDistance
		geometryOp.MinDistance = minDistance
		cmd.Location = NearOp{
			Near: geometryOp,
		}
	case OP_TYPE_NEAR_SPHERE:
		geometryOp.MaxDistance = maxDistance
		geometryOp.MinDistance = minDistance
		return NearSphereOP{
			NearSphere: geometryOp,
		}, nil
	default:
		return nil, fmt.Errorf("op type %v not supported", opType)
	}

	return &cmd, nil
}
