// app/models/document.go
package models

import (
	"os"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/database/orm"


)
	
type Document struct {
	orm.Model
	Filename string
	Path     string
}


// Method: Full path on disk
func (d *Document) FullPath() string {
	return facades.Storage().Path(d.Path)
}

// Method: Check if file exists
func (d *Document) Exists() bool {
	_, err := os.Stat(d.FullPath())
	return !os.IsNotExist(err)
}


var DocumentErrorMessage = map[string]string{
	"not_found":         "Document not found.",
	"create_failed":     "Failed to create the document.",
	"validation_failed": "Invalid input. Please check the fields and try again.",
	"invalid_request":   "Invalid request body. Please check your JSON format.",
	"internal_error":    "Something went wrong. Please try again later.",
}