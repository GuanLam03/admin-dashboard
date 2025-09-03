package permissions

type Permission struct {
	Label  string
	Object string
	Action string
}

var Permissions = map[string]Permission{
	"roles.index":  {Label: "View Roles", Object: "/roles", Action: "GET"},
	"roles.create": {Label: "Create Role", Object: "/roles", Action: "POST"},
	"roles.edit":   {Label: "Edit Role", Object: "/roles/:id", Action: "POST"},
	// "roles.delete": {Label: "Delete Role", Object: "/roles/:id", Action: "DELETE"},
}

// Helper function for backend
func PermissionKeyToObjectAction(key string) (string, string) {
	if perm, ok := Permissions[key]; ok {
		return perm.Object, perm.Action
	}
	return "", ""
}

// Optional: helper for frontend to get label
func PermissionKeyToLabel(key string) string {
	if perm, ok := Permissions[key]; ok {
		return perm.Label
	}
	return ""
}

// reverse lookup: Object + Action → Key
func PermissionObjectActionToKey(object, action string) string {
	for key, perm := range Permissions {
		if perm.Object == object && perm.Action == action {
			return key
		}
	}
	return ""
}
