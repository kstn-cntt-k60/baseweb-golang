package main

func AccountRoutes(root *Root) {
	root.PostAuthorized(
		"/api/account/add-party",
		"VIEW_EDIT_PARTY",
		root.account.AddPartyHandler)

	root.GetAuthorized(
		"/api/account/view-person",
		"VIEW_EDIT_PARTY",
		root.account.ViewPersonHandler)

	root.GetAuthorized(
		"/api/account/view-customer",
		"VIEW_EDIT_PARTY",
		root.account.ViewCustomerHandler)

	root.PostAuthorized(
		"/api/account/update-person",
		"VIEW_EDIT_PARTY",
		root.account.UpdatePersonHandler)

	root.PostAuthorized(
		"/api/account/delete-person",
		"VIEW_EDIT_PARTY",
		root.account.DeletePersonHandler)

	root.PostAuthorized(
		"/api/account/update-customer",
		"VIEW_EDIT_PARTY",
		root.account.UpdateCustomerHandler)

	root.PostAuthorized(
		"/api/account/delete-customer",
		"VIEW_EDIT_PARTY",
		root.account.DeleteCustomerHandler)

	root.GetAuthorized(
		"/api/account/query-simple-person",
		"VIEW_EDIT_PARTY",
		root.account.QuerySimplePersonHandler)

	root.PostAuthorized(
		"/api/account/add-user-login",
		"VIEW_EDIT_PARTY",
		root.account.AddUserLogin)

	root.GetAuthorized(
		"/api/account/view-user-login",
		"VIEW_EDIT_PARTY",
		root.account.ViewUserLoginHandler)

	root.PostAuthorized(
		"/api/account/update-user-login",
		"VIEW_EDIT_PARTY",
		root.account.UpdateUserLoginHandler)

	root.PostAuthorized(
		"/api/account/delete-user-login",
		"VIEW_EDIT_PARTY",
		root.account.DeleteUserLoginHandler)
}
