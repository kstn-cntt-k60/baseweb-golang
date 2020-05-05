package main

func ExportRoutes(root *Root) {
	root.PostAuthorized(
		"/api/export/export-sale-order-item",
		"EXPORT",
		root.export.ExportSaleOrderItemHandler,
	)

	root.GetAuthorized(
		"/api/export/view-exportable-sales-order",
		"EXPORT",
		root.export.ViewExportableSalesOrderHandler,
	)

	root.GetAuthorized(
		"/api/export/view-completed-sales-order",
		"EXPORT",
		root.export.ViewCompletedSalesOrderHandler,
	)

	root.GetAuthorized(
		"/api/export/view-single-sale-order",
		"EXPORT",
		root.order.ViewSingleSaleOrderHandler,
	)

	root.PostAuthorized(
		"/api/export/complete-sales-order",
		"EXPORT",
		root.export.CompleteSalesOrderHandler,
	)
}
