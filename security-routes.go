package main

func SecurityRoutes(root *Root) {
	root.PostAuthenticated("/api/login", root.security.LoginHandler)

	root.GetAuthorized(
		"/api/security/permission",
		"VIEW_EDIT_SECURITY_PERMISSION",
		root.security.SecurityPermissionHandler)

	root.PostAuthorized(
		"/api/security/save-group-permissions",
		"VIEW_EDIT_SECURITY_PERMISSION",
		root.security.SaveGroupPermissonsHandler)

	root.PostAuthorized(
		"/api/security/add-security-group",
		"VIEW_EDIT_SECURITY_GROUP",
		root.security.AddSecurityGroupHandler)

	root.GetAuthorized(
		"/api/security/user-login-info/{id}",
		"VIEW_EDIT_SECURITY_GROUP",
		root.security.UserLoginInfoHandler)

	root.PostAuthorized(
		"/api/security/save-user-login-security-groups",
		"VIEW_EDIT_SECURITY_GROUP",
		root.security.SaveUserLoginGroupsHandler)
}
