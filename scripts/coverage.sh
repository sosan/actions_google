go test -v ./pkg/... ./tests/... \
  -coverpkg=./pkg/... \
  -coverprofile=coverage.out \
  -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
echo "Report generated in coverage.html"