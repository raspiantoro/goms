# Goms

Goms is a command-line utility designed to work seamlessly with Gormigrate and Gormseeder, providing an efficient and structured approach to creating and managing database migrations and seeders for projects using the Gorm ORM in Golang. While Gormigrate and Gormseeder handle the intricacies of database migrations and seeding, Goms focuses on managing and generating migration and seeder files.

With Goms, you can easily generate migration and seeder files that are organized and maintainable. Goms is particularly useful for developers who need a reliable and straightforward way to manage database migration and data population in their applications, complementing Gorm's powerful ORM capabilities with an equally robust database migration and seeding solution.

## How to install
```bash
go install github.com/raspiantoro/goms@latest
```

## How to use

Before adding your seed files, you need to initialize Gomser within your application.
```bash
goms init
```

Then, you can add a new seed with the following command:

```bash
goms add seed "insert product"
```

It will produce a new seed file in your seeds directory with the following content:

```go
// Code generated by Goms (gorms).

package seeds

import (
	"gorm.io/gorm"
)

// don't rename this function
func (s *Seeds) SeedInsertProduct(db *gorm.DB) error {
	// place your seed code here
	return nil
}

// don't rename this function
func (s *Seeds) RollbackInsertProduct(db *gorm.DB) error {
	// place your rollback code here
	return nil
}
```

You can update both your seed function and rollback function. Below is an example:
```go
// don't rename this function
func (s *Seeds) SeedInsertProduct(db *gorm.DB) error {
	type Product struct {
		gorm.Model
		ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
		Name string
	}

	products := []*Product{
		{
			ID:   uuid.MustParse("e5024ae0-c9e0-40f9-b2b7-9813e125cb16"),
			Name: "Amazon Fire TV Stick",
		},
		{
			ID:   uuid.MustParse("118ae13f-afd0-4433-89e3-bce9770c4cc9"),
			Name: "Samsung Galaxy Tab S6",
		},
	}

	for _, product := range products {
		result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&product)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// don't rename this function
func (s *Seeds) RollbackInsertProduct(db *gorm.DB) error {
	type Product struct {
		gorm.Model
		ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
		Name string
	}

	productIDs := []uuid.UUID{
		uuid.MustParse("e5024ae0-c9e0-40f9-b2b7-9813e125cb16"),
		uuid.MustParse("118ae13f-afd0-4433-89e3-bce9770c4cc9"),
	}

	result := db.Unscoped().Delete(&Product{}, productIDs)

	return result.Error
}
```