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
