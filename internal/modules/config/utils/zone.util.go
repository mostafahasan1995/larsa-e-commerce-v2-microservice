package config_utils

func IsPointInRing(point []float64, ring [][]float64) bool {
	longitude := point[0]
	latitude := point[1]

	inside := false
	numVertices := len(ring)
	if numVertices < 3 { // Need at least 3 vertices for a valid ring
		return false
	}

	// The last vertex connects to the first vertex to close the ring
	j := numVertices - 1
	for i := 0; i < numVertices; i++ {
		// Vertex i coordinates (longitude, latitude)
		xi, yi := ring[i][0], ring[i][1]
		// Vertex j coordinates (longitude, latitude)
		xj, yj := ring[j][0], ring[j][1]

		// Check if the point's latitude is between the edge's y-coordinates
		// And if the point is to the left of the edge (using the ray casting formula)
		// Add a small epsilon for floating point comparisons if needed, but usually handled by the inequality checks
		intersect := ((yi > latitude) != (yj > latitude)) &&
			(longitude < (xj-xi)*(latitude-yi)/(yj-yi)+xi)

		if intersect {
			inside = !inside // Flip the inside status
		}

		// Move to the next edge
		j = i
	}
	return inside
}

func IsPointInPolygon(latitude, longitude float64, polygonData [][][]float64) bool {
	point := []float64{longitude, latitude}

	if len(polygonData) == 0 {
		return false
	}

	// Check if the point is in any of the rings
	for _, ring := range polygonData {
		if IsPointInRing(point, ring) {
			return true
		}
	}

	return false
}
