package main

func OrderRoutes(root *Root) {
	root.GetAuthorized(
		"/api/order/view-customer-store-by-customer",
		"VIEW_EDIT_ORDER",
		root.order.ViewCustomerStoreByCustomerHandler)

	root.PostAuthorized(
		"/api/order/add-order",
		"VIEW_EDIT_ORDER",
		root.order.AddOrderHandler)

	root.GetAuthorized(
		"/api/order/view-product-info-by-warehouse",
		"VIEW_EDIT_ORDER",
		root.order.ViewProductInfoByWarehouseHandler)

	root.GetAuthorized(
		"/api/order/view-sale-order",
		"VIEW_EDIT_ORDER",
		root.order.ViewSaleOrderHandler)

	root.GetAuthorized(
		"/api/order/view-single-sale-order",
		"VIEW_EDIT_ORDER",
		root.order.ViewSingleSaleOrderHandler)
}
