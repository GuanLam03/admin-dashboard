package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250825133209CreateUserRolesTable struct{}

func (r *M20250825133209CreateUserRolesTable) Signature() string {
	return "20250825133209_create_user_roles_table"
}

func (r *M20250825133209CreateUserRolesTable) Up() error {
	if !facades.Schema().HasTable("user_roles") {
		return facades.Schema().Create("user_roles", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id")
			table.UnsignedBigInteger("role_id")
			table.TimestampsTz()

			table.Index("user_id")
			table.Index("role_id")
		})
	}
	return nil
}

func (r *M20250825133209CreateUserRolesTable) Down() error {
	return facades.Schema().DropIfExists("user_roles")
}
