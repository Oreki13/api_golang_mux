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

# Screenshot

1. Get Methode 

<p align="center">
    <span>
      <img src="https://raw.githubusercontent.com/Oreki13/api_golang_mux/master/ss/get.png" width="500px" />
      &nbsp;&nbsp;
    </span>
  </p>

2. Post Methode

<p align="center">
    <span>
      <img src="https://raw.githubusercontent.com/Oreki13/api_golang_mux/master/ss/post.png" width="500px" />
      &nbsp;&nbsp;
    </span>
  </p>

3. PUT Methode

<p align="center">
    <span>
      <img src="https://raw.githubusercontent.com/Oreki13/api_golang_mux/master/ss/edit.png" width="500px" />
      &nbsp;&nbsp;
    </span>
  </p>

4. Delete Methode

<p align="center">
    <span>
      <img src="https://raw.githubusercontent.com/Oreki13/api_golang_mux/master/ss/delete.png" width="500px" />
      &nbsp;&nbsp;
    </span>
  </p>

