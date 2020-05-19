package main

func SalesmanRoutes(root *Root) {
	root.GetAuthorized(
		"/api/salesman/view-schedule",
		"SALESMAN_CHECKIN",
		root.salesman.ViewScheduleHandler)

	root.GetAuthorized(
		"/api/salesman/view-checkin-history",
		"SALESMAN_CHECKIN",
		root.salesman.ViewCheckinHistoryHandler)

	root.PostAuthorized(
		"/api/salesman/add-checkin",
		"SALESMAN_CHECKIN",
		root.salesman.AddCheckinHandler)
}
