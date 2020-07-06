package main

func ScheduleRoutes(root *Root) {
	root.GetAuthorized(
		"/api/schedule/view-store-city",
		"VIEW_EDIT_SALESMAN",
		root.schedule.ViewStoreCityHandler)
}
