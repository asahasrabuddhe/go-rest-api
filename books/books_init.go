package books

import (
    "github.com/asahasrabuddhe/rest-api/server"
)

func init() {
    server.Mount("/books", Resource{}.Routes())
}
