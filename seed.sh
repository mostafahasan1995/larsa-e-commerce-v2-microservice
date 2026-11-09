# Build the seed command
go build -o seed cmd/seed/main.go

# Run the seed command with the required flag
./seed -onlySystemData

rm ./seed