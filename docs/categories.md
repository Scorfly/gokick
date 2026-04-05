## Get Categories

`GetCategories` calls Kick’s **`GET /public/v2/categories`** endpoint. Results are cursor-paginated (`Pagination.NextCursor`). Filters are optional: `AddName`, `AddTag`, `AddID`, `SetLimit` (1–1000, API default 25), and `SetCursor` for the next page.

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	categories, err := client.GetCategories(context.Background(), gokick.NewCategoryListFilter().AddName("Gaming Cafe Simulator").SetLimit(25))
	if err != nil {
		log.Fatalf("Failed to fetch categories: %v", err)
	}

	spew.Dump("categories", categories)
```

### Pagination (cursor)

```go
	client, _ := gokick.NewClient(&gokick.ClientOptions{
		UserAccessToken: "xxxx",
	})

	first, err := client.GetCategories(context.Background(), gokick.NewCategoryListFilter().SetLimit(10))
	if err != nil {
		log.Fatalf("Failed to fetch categories: %v", err)
	}

	if first.Pagination.NextCursor != "" {
		next, err := client.GetCategories(context.Background(), gokick.NewCategoryListFilter().SetCursor(first.Pagination.NextCursor).SetLimit(10))
		if err != nil {
			log.Fatalf("Failed to fetch next page: %v", err)
		}
		spew.Dump("next page", next)
	}
```

### Filters by tag or id

```go
	filter := gokick.NewCategoryListFilter().AddTag("fps").AddID(9569)
	categories, err := client.GetCategories(context.Background(), filter)
```

output (shape may vary with live API data)

```
(string) (len=10) "categories"
(gokick.CategoriesResponseWrapper) {
 Result: ([]gokick.CategoryResponse) (len=2 cap=2) {
  (gokick.CategoryResponse) {
   ID: (int) 9569,
   Name: (string) (len=21) "Gaming Cafe Simulator",
   Tags: ([]string) (len=1 cap=1) {
    (string) (len=4) "simu"
   },
   Thumbnail: (string) (len=115) "https://files.kick.com/images/subcategories/9569/banner/conversion/db0164a6-05e9-42cd-8066-1c30a0d999a6-banner.webp"
  },
  (gokick.CategoryResponse) {
   ID: (int) 4688,
   Name: (string) (len=19) "Miniature Wargaming",
   Tags: ([]string) <nil>,
   Thumbnail: (string) (len=115) "https://files.kick.com/images/subcategories/4688/banner/conversion/9042d2d7-8d2b-4550-ad00-21d4420b627f-banner.webp"
  }
 },
 Pagination: (gokick.Pagination) {
  NextCursor: (string) (len=12) "eyJjX2N1cnNvciI6MTB9"
 }
}
```

## Get Category

`GetCategory` resolves a single category by id (it lists via v2 with `AddID` and returns the first row).

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
  Tags: ([]string) (len=1 cap=1) {
   (string) (len=4) "simu"
  },
  Thumbnail: (string) (len=115) "https://files.kick.com/images/subcategories/9569/banner/conversion/db0164a6-05e9-42cd-8066-1c30a0d999a6-banner.webp"
 }
}
```
