package util

import (
	"errors"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/nyaruka/phonenumbers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Identifiable interface {
	GetId() string
}

func Paginate(r *http.Request) (sk int64, lim int64, er error) {
	page := r.URL.Query().Get("page")
	perPage := r.URL.Query().Get("perPage")
	var skip, limit int64
	var convErr error
	if perPage == "" {
		limit = 40
	} else {
		limit, convErr = strconv.ParseInt(perPage, 10, 64)
		if convErr != nil {
			return 0, 0, convErr
		}
	}
	if page == "" {
		skip = 0
	} else {
		skip, convErr = strconv.ParseInt(page, 10, 64)
		if convErr != nil {
			return 0, 0, convErr
		}
		skip = (skip - 1) * limit
	}

	return skip, limit, nil
}

func PaginatePage(r *http.Request) (int, int) {
	perPage := r.URL.Query().Get("perPage")
	if perPage == "" {
		perPage = "10"
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		return 0, 0
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return 0, 0
	}
	return perPageInt, pageInt
}

func CacheKeyForMulti(key enums.CollectionName, db, path, rawQuery string) string {
	pattern := "%s:p:c:%s:%s?%s"
	return fmt.Sprintf(pattern, key, db, path, rawQuery)
}

func CacheKeyForOne(key enums.CollectionName, db, id string) string {
	pattern := "%s:g:c:%s:%s"
	return fmt.Sprintf(pattern, key, db, id)
}

func CacheKeyDeleteMul(key enums.CollectionName, db string) string {
	pattern := "%s:p:c:%s:*"
	return fmt.Sprintf(pattern, key, db)
}
func CacheKeyDeleteOne(key enums.CollectionName, db, id string) string {
	pattern := "%s:g:c:%s:%s"
	return fmt.Sprintf(pattern, key, db, id)
}

func InSlice[T comparable](ele T, s []T) bool {
	for _, e := range s {
		if e == ele {
			return true
		}
	}

	return false

}

// GetCurrentMonthDates returns the start and end dates of the current month in RFC3339 format
func GetCurrentMonthDates() (string, string) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return startOfMonth.Format(time.RFC3339), endOfMonth.Format(time.RFC3339)
}

// GetDateOrDefault returns the provided date or defaults to current month dates
func GetDateOrDefault(date string, isStartDate bool) string {
	if date != "" {
		return date
	}

	startDate, endDate := GetCurrentMonthDates()
	if isStartDate {
		return startDate
	}
	return endDate
}

func GetBsonFromString(query string) (bson.M, error) {
	var bs bson.M = bson.M{}
	if query != "" {
		err := bson.UnmarshalExtJSON([]byte(query), true, &bs)
		if err != nil {
			return nil, err
		}
	}

	return bs, nil
}

func GeneratePassword(passWordLength, minSpecialChar, minUpper, minNum int) string {
	var (
		lowerCharSet   = "abcdedfghijklmnopqrst"
		upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialCharSet = "!@#$%&*"
		numberSet      = "0123456789"
		allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
	)
	var password strings.Builder
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < minSpecialChar; i++ {
		index := r.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[index]))
	}

	for i := 0; i < minUpper; i++ {
		index := r.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[index]))
	}

	for i := 0; i < minNum; i++ {
		index := r.Intn(len(numberSet))
		password.WriteString(string(numberSet[index]))
	}

	remainLength := passWordLength - minSpecialChar - minUpper - minNum

	for i := 0; i < remainLength; i++ {
		index := r.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[index]))
	}

	return password.String()

}

// as is date or time
func GetDateAs(date time.Time, as string) time.Time {
	if as == "date" {
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	} else if as == "time" {
		return time.Date(0, 0, 0, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), time.UTC)
	}

	return date

}

type ReqArgs struct {
	Sort    bson.M
	GroupBy bson.M
	//any other args needed
}

func (a *ReqArgs) SetSort(r *http.Request) error {
	sortBy := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")

	if sortBy != "" && sortOrder != "" {
		sort := bson.M{}
		order, convErr := strconv.ParseInt(sortOrder, 10, 64)
		if convErr != nil {
			return convErr
		}
		sort[sortBy] = order
		a.Sort = sort
	}

	return nil
}

func (a *ReqArgs) SetGroupBy(r *http.Request) error {
	groupBy := r.URL.Query().Get("groupBy")

	if groupBy != "" {
		group := bson.M{
			"_id":     "$" + groupBy,
			"records": bson.M{"$push": "$$ROOT"},
		}

		a.GroupBy = group
	}

	return nil

}

func GetWeekdayNames(startDate, endDate time.Time) []string {
	var weekdays []string

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		weekdays = append(weekdays, d.Weekday().String())
	}

	return weekdays
}

func BsonAToSlice[T any](bsonArray bson.A) ([]T, error) {
	var result []T
	for _, ele := range bsonArray {
		t, ok := ele.(T)
		if !ok {
			return nil, errors.New("error type")
		}

		result = append(result, t)
	}
	return result, nil
}

func Contains[T comparable](arr []T, key T) bool {
	for _, value := range arr {
		if value == key {
			return true
		}
	}
	return false
}

func SliceFilter[T any](arr []T, callBack func(el T) bool) []T {
	var result []T
	for _, e := range arr {
		ok := callBack(e)
		if ok {
			result = append(result, e)
		}
	}

	return result
}

func SliceFind[T any](arr []T, callback func(el T) bool) (T, bool) {
	var result T
	var ok bool
	for _, e := range arr {
		ok = callback(e)
		if ok {
			result = e
			break
		}
	}

	return result, ok
}

func RoundToDecimalPlaces(value float64, decimalPlaces int) float64 {
	factor := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*factor) / factor
}

func IsValidEmail(email string) bool {
	// Define the regex pattern for a valid email address
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern
	re := regexp.MustCompile(emailPattern)

	// Match the email address against the pattern
	return re.MatchString(email)
}

func IsValidPhoneNumber(phone string) bool {
	num, err := phonenumbers.Parse(phone, "") // empty region since number includes country code
	if err != nil {
		return false
	}

	if !phonenumbers.IsValidNumber(num) {
		return false
	}

	return true
}

func GenerateUniqueString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func HasDuplicateIDs[T Identifiable](slice []T) bool {
	seen := make(map[string]bool)
	for _, item := range slice {
		if seen[item.GetId()] {
			return true
		}
		seen[item.GetId()] = true
	}
	return false
}

func SliceReorder[T Identifiable](s []T, oldIndex, newIndex int) ([]T, error) {
	if oldIndex >= 0 && oldIndex < len(s) && newIndex >= 0 && newIndex < len(s) {
		// Remove the element from the old position
		ele := s[oldIndex]
		s = append(s[:oldIndex], s[oldIndex+1:]...)
		// Insert the element at the new position
		s = append(s[:newIndex], append([]T{ele}, s[newIndex:]...)...)
		return s, nil
	} else {
		return s, errors.New("invalid indices")
	}
}

func ToObjectIds(ids []string) []primitive.ObjectID {
	objectIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objectID, _ := primitive.ObjectIDFromHex(id)
		objectIDs[i] = objectID
	}
	return objectIDs
}

func ExtractUniqueValues[T comparable](slice []T) []T {
	uniqueValues := make(map[T]bool)
	for _, value := range slice {
		uniqueValues[value] = true
	}

	uniqueSlice := make([]T, 0, len(uniqueValues))
	for value := range uniqueValues {
		uniqueSlice = append(uniqueSlice, value)
	}

	return uniqueSlice
}

// HasDuplicateValues checks if there are duplicate values in a slice
func HasDuplicateValues[T comparable](slice []T) bool {
	seen := make(map[T]bool)
	for _, value := range slice {
		if seen[value] {
			return true
		}
		seen[value] = true
	}
	return false
}

// Filter filters a slice based on a callback function
func Filter[T any](slice []T, callback func(el T, index int) bool) []T {
	var result []T
	for index, e := range slice {
		ok := callback(e, index)
		if ok {
			result = append(result, e)
		}
	}
	return result
}

func Map[T any, R any](slice []T, callback func(el T) R) []R {
	var result []R
	for _, e := range slice {
		result = append(result, callback(e))
	}
	return result
}

func Find[T any](slice []T, callback func(el T) bool) *T {
	for _, e := range slice {
		ok := callback(e)
		if ok {
			return &e
		}
	}
	return nil
}

// IsNumber checks if a value is a number by attempting to convert it to float64
func IsNumber(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true

	default:
		return false
	}
}

func IsArray(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Array || reflect.TypeOf(v).Kind() == reflect.Slice
}

func IsBoolean(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Bool
}

func IsString(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.String
}

func IsObject(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}

func IsEnum(value string, enumType reflect.Type) bool {
	// Create a zero value of the enum type
	zeroValue := reflect.Zero(enumType).Interface()

	// Try to convert the string to the enum type
	enumValue := reflect.ValueOf(zeroValue).Convert(reflect.TypeOf(value)).Interface()

	// Check if the conversion was successful
	return enumValue == value
}

func GenerateCode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return time.Now().Format("20060102150405") + "_" + string(b)
}

func ToInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("unexpected type %T", val)
	}
}
