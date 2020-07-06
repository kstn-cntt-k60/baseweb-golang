package main

func FacilityRoutes(root *Root) {
	root.PostAuthorized(
		"/api/facility/add-warehouse",
		"VIEW_EDIT_FACILITY",
		root.facility.AddWarehouseHandler)

	root.GetAuthorized(
		"/api/facility/view-warehouse",
		"VIEW_EDIT_FACILITY",
		root.facility.ViewWarehouseHandler)

	root.GetAuthorized(
		"/api/facility/get-warehouse/{warehouseId}",
		"VIEW_EDIT_FACILITY",
		root.facility.GetWarehouseHandler)

	root.PostAuthorized(
		"/api/facility/update-warehouse",
		"VIEW_EDIT_FACILITY",
		root.facility.UpdateWarehouseHandler)

	root.PostAuthorized(
		"/api/facility/delete-warehouse",
		"VIEW_EDIT_FACILITY",
		root.facility.DeleteWarehouseHandler)

	root.GetAuthorized(
		"/api/facility/view-customer-store",
		"VIEW_EDIT_FACILITY",
		root.facility.ViewCustomerStoreHandler)

	root.GetAuthorized(
		"/api/facility/query-simple-customer",
		"VIEW_EDIT_FACILITY",
		root.facility.QuerySimpleCustomerHandler)

	root.PostAuthorized(
		"/api/facility/add-customer-store",
		"VIEW_EDIT_FACILITY",
		root.facility.AddCustomerStoreHandler)

	root.PostAuthorized(
		"/api/facility/update-customer-store",
		"VIEW_EDIT_FACILITY",
		root.facility.UpdateCustomerStoreHandler)

	root.PostAuthorized(
		"/api/facility/delete-customer-store",
		"VIEW_EDIT_FACILITY",
		root.facility.DeleteCustomerStoreHandler)

	root.GetAuthorized(
		"/api/facility/view-all-customer-store",
		"VIEW_EDIT_FACILITY",
		root.facility.ViewAllCustomerStoreHandler)
}
