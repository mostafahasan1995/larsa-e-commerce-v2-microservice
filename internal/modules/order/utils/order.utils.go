package order_utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateOrderCode() string {
	// Generate a unique 6-digit numeric order code
	objID := primitive.NewObjectID()
	// Combine timestamp (last 3 digits) with random bytes for uniqueness
	timestamp := objID.Timestamp().Unix() % 1000
	randomBytes := int64(0)
	for i := 6; i < 12; i++ {
		randomBytes = (randomBytes*10 + int64(objID[i]%10)) % 1000
	}
	// Combine timestamp and random bytes to create unique 6-digit number
	uniqueNum := (timestamp*1000 + randomBytes) % 1000000
	return "ORD-" + fmt.Sprintf("%06d", uniqueNum)
}
