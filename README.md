# Golang-products-API
A simple REST API for a catalog of products

### This project is based on Echo framework for Go language
### To see available endpoints, launch application and visit the "/" path

#### HTTP
- `github.com/labstack/echo/v4` - used for developing a REST service;

#### Validation
- `gopkg.in/go-playground/validator.v10` - for _data validation_  

#### Database
- `database/sql`
- `github.com/mattn/go-sqlite3` - sqlite3 as SQL driver

#### Configuration
- `flag` - to set app launch flags
- `github.com/jinzhu/configor` - config load from yaml file

#### Logging
- `github.com/sirupsen/logrus` - app log

#### Testing
- `testing` - testing
- `github.com/golang/mock/gomock` - mock generation
- `github.com/stretchr/testify/assert` - _assert_ in testing

### App launch
- `go run main.go -conf="Path-to-conf-file"` - start the main application | _conf_ flag is used to point to config file. By default "./config/config.yaml" is used
- `go test -v ./test/` - Performs testing
