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

// GeoIntersectsOp selects documents whose geospatial data intersects with a specified GeoJSON object; 
// i.e. where the intersection of the data and the specified object is non-empty.
type GeoIntersectsOp struct {
	GeoIntersects interface{} `bson:"$geoIntersects"`
}

// GeoWithinOP selects documents with geospatial data that exists entirely within a specified shape.
// The specified shape can be either a GeoJSON Polygon (either single-ringed or multi-ringed), 
// a GeoJSON MultiPolygon, or a shape defined by legacy coordinate pairs. 
// The $geoWithin operator uses the $geometry operator to specify the GeoJSON object.
type GeoWithinOP struct {
	GeoWithin interface{} `bson:"$geoWithin"`
}

// NearOp specifies a point for which a geospatial query returns the documents from nearest to farthest. 
// The $near operator can specify either a GeoJSON point or legacy coordinate point.
type NearOp struct {
	Near interface{} `bson:"$near"`
}

// NearSphereOP specifies a point for which a geospatial query returns the documents from nearest to farthest. 
// MongoDB calculates distances for $nearSphere using spherical geometry.
type NearSphereOP struct {
	NearSphere interface{} `bson:"$nearSphere"`
}

type MaxDistanceOp struct {
	MaxDistance int `bson:"$maxDistance"`
}

type MinDistanceOp struct {
	MinDistance int `bson:"$minDistance"`
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
