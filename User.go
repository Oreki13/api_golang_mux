package restapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func (c *InitAPI) ListUser(ctx context.Context, req *GetUsers) (*GetUsers, error) {
	limit := 10

	if req.Limit != 0 {
		limit = int(req.Limit)
	}

	rows, err := c.Db.Query(`
		SELECT id, 
			username, 
			email,
			status, 
			role_id,
			created_at,
			updated_at
		FROM users LIMIT $1
	`, limit)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var items []*User
	for rows.Next() {
		var item User
		var updateTime sql.NullString
		var status string
		err = rows.Scan(&item.Id,
			&item.Username,
			&item.Email,
			&status,
			&item.RoleId,
			&item.CreatedAt,
			&updateTime,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		item.UpdatedAt = updateTime.String

		items = append(items, &item)
	}

	if len(items) == 0 {
		return nil, errors.New("user-not-found")
	}

	return &GetUsers{
		List: items,
	}, nil
}

// CreateUser for creating user
func (c *InitAPI) CreateUser(ctx context.Context, req *User, rolesId string) (*UserId, error) {
	var id string
	roles, err := c.GetRoles(rolesId)
	if err != nil {
		log.Println(err)
		if err.Error() == "no rows in result set" {
			return nil, errors.New("ERROR-NO-ADMIN-FOUND")
		}
		return nil, err
	}

	if roles != "ADMIN" {
		return nil, errors.New("invalid-roles")
	}

	status := strconv.Itoa(req.Status)
	err = c.Db.QueryRow(`INSERT INTO users (username, email, status, role_id, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Username, req.Email, status, rolesId, time.Now()).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &UserId{
		Id: id,
	}, nil
}

func (c *InitAPI) EditUser(ctx context.Context, req *UserEdit, rolesId string) (*UserEdit, error) {
	// var id string
	roles, err := c.GetRoles(rolesId)

	if len(req.Id) < 1 {
		return nil, errors.New("ID-NOT-FOUND")
	}

	IDUser, errUser := c.GetUserId(req.Id)

	if err != nil {
		log.Println(err)
		if err.Error() == "no rows in result set" {
			return nil, errors.New("ERROR-NO-ADMIN-FOUND")
		}
		return nil, err
	}

	if errUser != nil {
		log.Println(errUser)
		if errUser.Error() == "no rows in result set" {
			return nil, errors.New("ERROR-NO-USER-FOUND")
		}
		return nil, errUser
	}

	if roles != "ADMIN" {
		return nil, errors.New("invalid-roles")
	}

	if len(req.Username) < 1 {
		req.Username = IDUser[0].Username
	}

	if len(req.Email) < 1 {
		req.Email = IDUser[0].Email
	}

	if len(req.Status) < 1 {
		req.Status = IDUser[0].Status
	}
	if len(req.RoleId) < 1 {
		req.RoleId = IDUser[0].RoleId
	}

	_, err = c.Db.Exec(`UPDATE users SET username = $1, email = $2, status = $3, role_id = $4, updated_at = $5 WHERE id = $6`,
		req.Username, req.Email, req.Status, req.RoleId, time.Now(), req.Id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &UserEdit{
		Id:       req.Id,
		Username: req.Username,
		Email:    req.Email,
		Status:   req.Status,
		RoleId:   req.RoleId,
	}, nil
}

func (c *InitAPI) DeleteUser(ctx context.Context, req *UserId, rolesId string) (*UserId, error) {
	// var id string
	roles, err := c.GetRoles(rolesId)

	if len(req.Id) < 1 {
		return nil, errors.New("ID-NOT-FOUND")
	}

	_, errUser := c.GetUserId(req.Id)

	if err != nil {
		log.Println(err)
		if err.Error() == "no rows in result set" {
			return nil, errors.New("ERROR-NO-ADMIN-FOUND")
		}
		return nil, err
	}

	if errUser != nil {
		log.Println(errUser)
		if errUser.Error() == "no rows in result set" {
			return nil, errors.New("ERROR-NO-USER-FOUND")
		}
		return nil, errUser
	}

	if roles != "ADMIN" {
		return nil, errors.New("invalid-roles")
	}

	_, err = c.Db.Exec(`DELETE FROM users WHERE id = $1`, req.Id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &UserId{
		Id: req.Id,
	}, nil
}

func (c *InitAPI) GetRoles(id string) (string, error) {
	var roles string
	err := c.Db.QueryRow(`SELECT roles FROM roles WHERE id = $1`, id).Scan(&roles)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return roles, nil
}

func (c *InitAPI) GetUserId(id string) ([]UserEdit, error) {
	var myUser UserEdit
	var user []UserEdit

	rows, err := c.Db.Query(`SELECT username, email, status, role_id FROM users WHERE id = $1`, id)

	if err != nil {
		log.Println(err)
		return user, err
	}

	for rows.Next() {
		if err := rows.Scan(&myUser.Username, &myUser.Email, &myUser.Status, &myUser.RoleId); err != nil {
			log.Fatal(err.Error())
		} else {
			user = append(user, myUser)
		}

	}

	return user, nil
}

func (c *InitAPI) InsertProfilePhoto(ctx context.Context, req *FileItem) (*UserId, error) {
	if !c.GetCustomerById(req.UserId) {
		return nil, errors.New("user-not-found")
	}

	if _, _, err := c.GetProfilePhotoById(req.UserId); err == nil {
		return nil, errors.New("ERROR-USER-UDAH-ADA-FOTO")
	}

	var profileId string
	err := c.Db.QueryRow(`INSERT INTO profile_photo (user_id, filename, file_type, size) VALUES ($1, $2, $3, $4) RETURNING id`,
		req.UserId, req.Filename, req.FileType, req.FileSize,
	).Scan(&profileId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	filename := fmt.Sprintf("assert/%s", req.Filename)

	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer file.Close()
	_, err = io.Copy(file, req.File)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &UserId{
		Id: profileId,
	}, nil
}

func (c *InitAPI) GetProfilePhotoById(id string) (string, string, error) {
	var filename, fileType string
	err := c.Db.QueryRow(`SELECT filename, file_type FROM profile_photo WHERE user_id = $1`, id).Scan(&filename, &fileType)

	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return filename, fileType, nil
}

func (c *InitAPI) GetProfilePhoto(ctx context.Context, req *GetFile) (io.Reader, string, error) {
	filename, fileType, err := c.GetProfilePhotoById(req.UserId)
	if err != nil {
		return nil, "", nil
	}
	url := fmt.Sprintf("assert/%s", filename)
	file, err := os.Open(url)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	return file, fileType, nil
}

func (c *InitAPI) GetCustomerById(id string) bool {
	var userId string
	err := c.Db.QueryRow(`SELECT username FROM users WHERE id = $1`, id).Scan(&userId)
	if err != nil {
		return false
	}

	return userId != ""
}
