package main

func ImportRoutes(root *Root) {
	root.GetAuthorized(
		"/api/import/view-product-by-warehouse",
		"IMPORT",
		root.importProduct.ViewProductByWarehouseHandler)

	root.PostAuthorized(
		"/api/import/add-inventory-item",
		"IMPORT",
		root.importProduct.AddInventoryItemHandler)

	root.GetAuthorized(
		"/api/import/view-inventory-by-warehouse",
		"IMPORT",
		root.importProduct.ViewInventoryByWarehouseHandler)

	root.GetAuthorized(
		"/api/import/view-inventory-by-product",
		"IMPORT",
		root.importProduct.ViewInventoryByProductHandler)
}
