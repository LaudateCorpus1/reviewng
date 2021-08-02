package db

// MIT License
//
// Copyright (c) 2021 Damian Zaremba
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

type User struct {
	Id          int
	Username    string
	Approved    bool
	Admin       bool
	LegacyCount int
}

func (db *Db) CreateUser(username string, approved bool, admin bool) (bool, error) {
	insert, err := db.db.Query("INSERT INTO users (username, admin, approved) VALUES (?, ?, ?)", username, admin, approved)
	if err != nil {
		return false, err
	}

	if err := insert.Close(); err != nil {
		return false, err
	}
	return true, nil
}

func (db *Db) LookupUserByName(username string) (*User, error) {
	results, err := db.db.Query(`SELECT id, admin, approved, legacy_count FROM users WHERE username = ?`, username)
	if err != nil {
		return nil, err
	}

	if !results.Next() {
		return nil, nil
	}

	user := User{Username: username}
	if err := results.Scan(&user.Id, &user.Admin, &user.Approved, &user.LegacyCount); err != nil {
		return nil, err
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Db) LookupUserById(id int) (*User, error) {
	results, err := db.db.Query(`SELECT username, admin, approved, legacy_count FROM users WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}

	if !results.Next() {
		return nil, nil
	}

	user := User{Id: id}
	if err := results.Scan(&user.Username, &user.Admin, &user.Approved, &user.LegacyCount); err != nil {
		return nil, err
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Db) FetchAllUsers() ([]*User, error) {
	results, err := db.db.Query("SELECT id, username, admin, approved, legacy_count FROM users")
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for results.Next() {
		user := &User{}
		if err := results.Scan(&user.Id, &user.Username, &user.Admin, &user.Approved, &user.LegacyCount); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := results.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *Db) CalculateTotalUserClassifications(user *User) (int, error) {
	results, err := db.db.Query("SELECT COUNT(*) FROM user_classification WHERE user_id = ?", user.Id)
	if err != nil {
		return -1, err
	}

	if !results.Next() {
		return -1, nil
	}

	var total int
	if err := results.Scan(&total); err != nil {
		return -1, err
	}

	if err := results.Close(); err != nil {
		return -1, err
	}

	return user.LegacyCount + total, nil
}