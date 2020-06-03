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
		"/api/sales-route/view-salesroute-config",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.ViewConfigHandler)

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

	root.PostAuthorized(
		"/api/sales-route/add-salesroute-config",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.AddConfigHandler)

	root.PostAuthorized(
		"/api/sales-route/update-salesroute-config",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.UpdateConfigHandler)

	root.PostAuthorized(
		"/api/sales-route/delete-salesroute-config",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.DeleteConfigHandler)

	root.PostAuthorized(
		"/api/schedule/add-schedule",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.AddScheduleHandler)

	root.GetAuthorized(
		"/api/schedule/view-schedule",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.ViewScheduleHandler)

	root.PostAuthorized(
		"/api/schedule/delete-schedule",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.DeleteScheduleHandler)

	root.GetAuthorized(
		"/api/schedule/get-schedule/{scheduleId}",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.GetScheduleHandler)

	root.GetAuthorized(
		"/api/schedule/view-clustering",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.ViewClusteringHandler)

	root.GetAuthorized(
		"/api/salesman/view-user-login",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.ViewUserLoginHandler)

	root.PostAuthorized(
		"/api/sales-route/delete-salesman",
		"VIEW_EDIT_SALESMAN",
		root.salesroute.DeleteSalesmanHandler)
}
