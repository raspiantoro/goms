package templates

func SeedTemplate() []byte {
	return []byte(`// Code generated by goms. DO NOT EDIT.

package {{ .SeederSubDir }}

import "github.com/raspiantoro/gormseeder"

var Seeds = []*gormseeder.Seed{}
`)
}

func SeederTemplate() []byte {
	return []byte(`// Code generated by goms. DO NOT RENAME THIS FILENAME.

package {{ .SeederSubDir }}

import (
	"gorm.io/gorm"
)

// don't rename this function
func {{ .UpFuncName }}(db *gorm.DB) error {
	// place your seed code here
	return nil
}

// don't rename this function
func {{ .DownFuncName }}(db *gorm.DB) error {
	// place your rollback code here
	return nil
}	
`)
}
