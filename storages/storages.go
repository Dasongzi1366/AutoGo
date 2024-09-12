package storages

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var database *sql.DB

func init() {
	dir, _ := os.Getwd()
	runDir := dir + "/"
	var err error
	database, err = sql.Open("sqlite3", runDir+"storages.db")
	if err == nil {
		_, err = os.Stat(runDir + "storages.db")
		if err != nil {
			_, _ = database.Exec(`CREATE TABLE IF NOT EXISTS storages ("key" TEXT NOT NULL PRIMARY KEY, "value" TEXT);`)
		}
	}
}

// Get 从本地存储中取出键值为key的数据并返回。
func Get(key string) string {
	var value string
	query := `SELECT value FROM storages WHERE key = ?`
	_ = database.QueryRow(query, key).Scan(&value)
	return value
}

// Put 把值value保存到本地存储中。
func Put(key, value string) {
	query := `INSERT OR REPLACE INTO storages (key, value) VALUES (?, ?)`
	_, _ = database.Exec(query, key, value)
}

// Remove 移除键值为key的数据。
func Remove(key string) {
	query := `DELETE FROM storages WHERE key = ?`
	_, _ = database.Exec(query, key)
}

// Contains 返回该本地存储是否包含键值为key的数据。
func Contains(key string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM storages WHERE key = ?)`
	_ = database.QueryRow(query, key).Scan(&exists)
	return exists
}

// Clear 移除该本地存储的所有数据。
func Clear() {
	query := `DELETE FROM storages`
	_, _ = database.Exec(query)
}
