package main

func SalesrouteRoutes(root *Root) {
	root.PostAuthorized(
		"/api/sales-route/add-salesman",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.AddSalesmanHandler)

	root.PostAuthorized(
		"/api/sales-route/add-planning-period",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.AddPlanningPeriodHandler)

	root.PostAuthorized(
		"/api/sales-route/update-planning-period",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.UpdatePlanningPeriodHandler)

	root.PostAuthorized(
		"/api/sales-route/delete-planning-period",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.DeletePlanningPeriodHandler)

	root.GetAuthorized(
		"/api/sales-route/view-salesman",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.ViewSalesmanHandler)

	root.GetAuthorized(
		"/api/sales-route/view-salesroute-planning",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.ViewPlanningPeriodHandler)

	root.GetAuthorized(
		"/api/sales-route/get-planning/{planningId}",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.GetPlanningPeriodHandler)
}
