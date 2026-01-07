## Get Categories

```go
    client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	categories, err := client.GetCategories(context.Background(), gokick.NewCategoryListFilter().SetQuery("gaming"))
	if err != nil {
		log.Fatalf("Failed to fetch categories: %v", err)
	}

	spew.Dump("categories", categories)
```

### Get Categories with pagination

```go
    client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	categories, err := client.GetCategories(context.Background(), gokick.NewCategoryListFilter().SetQuery("gaming").SetPage(2))
	if err != nil {
		log.Fatalf("Failed to fetch categories: %v", err)
	}

	spew.Dump("categories", categories)
```
output
```
(string) (len=10) "categories"
(gokick.CategoriesResponseWrapper) {
 Result: ([]gokick.CategoryResponse) (len=2 cap=2) {
  (gokick.CategoryResponse) {
   ID: (int) 9569,
   Name: (string) (len=21) "Gaming Cafe Simulator",
   Thumbnail: (string) (len=115) "https://files.kick.com/images/subcategories/9569/banner/conversion/db0164a6-05e9-42cd-8066-1c30a0d999a6-banner.webp"
  },
  (gokick.CategoryResponse) {
   ID: (int) 4688,
   Name: (string) (len=19) "Miniature Wargaming",
   Thumbnail: (string) (len=115) "https://files.kick.com/images/subcategories/4688/banner/conversion/9042d2d7-8d2b-4550-ad00-21d4420b627f-banner.webp"
  }
 }
}
```

## Get Category

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	category, _ := client.GetCategory(context.Background(), 9569)

	spew.Dump("category", category)
```
output
```
(string) (len=8) "category"
(gokick.CategoryResponseWrapper) {
 Result: (gokick.CategoryResponse) {
  ID: (int) 9569,
  Name: (string) (len=21) "Gaming Cafe Simulator",
  Thumbnail: (string) (len=115) "https://files.kick.com/images/subcategories/9569/banner/conversion/db0164a6-05e9-42cd-8066-1c30a0d999a6-banner.webp"
 }
}
```