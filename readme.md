# How To use

1. create server on postgres with this config 

```
	dbHost := "127.0.0.1"
	dbPass := "root"
	dbName := "postgres"
	dbPort := "5432"
	dbUser := "postgres"
```
2. migrate database on dir `migration` to your server

3. to create end point you must created on `route.go` and `StartHttp` func 

```
func StartHttp() http.Handler {
	api := createAPI()
	api.initDb()

	r := mux.NewRouter()
	r.HandleFunc("/api/user/list", api.HandleListUser).Methods("GET")
	r.HandleFunc("/api/user/create", api.HandleCreateUser).Methods("POST")

	return r
}
```

4. Create you logic func on `User.go` 
