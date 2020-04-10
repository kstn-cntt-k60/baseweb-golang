package main

func ProductRoutes(root *Root) {
	root.PostAuthorized(
		"/api/product/add-product",
		"VIEW_EDIT_PRODUCT",
		root.product.AddProductHandler)

	root.GetAuthorized(
		"/api/product/view-product",
		"VIEW_EDIT_PRODUCT",
		root.product.ViewProductHandler)

	root.PostAuthorized(
		"/api/product/update-product",
		"VIEW_EDIT_PRODUCT",
		root.product.UpdateProductHandler)

	root.PostAuthorized(
		"/api/product/delete-product",
		"VIEW_EDIT_PRODUCT",
		root.product.DeleteProductHandler)

	root.GetAuthorized(
		"/api/product/view-product-pricing",
		"VIEW_EDIT_PRODUCT",
		root.product.ViewProductPricingHandler)

	root.GetAuthorized(
		"/api/product/view-product-price",
		"VIEW_EDIT_PRODUCT",
		root.product.ViewProductPriceHandler)

	root.PostAuthorized(
		"/api/product/add-product-price",
		"VIEW_EDIT_PRODUCT",
		root.product.AddProductPriceHandler)
}
