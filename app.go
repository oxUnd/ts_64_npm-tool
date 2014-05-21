package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/martini-contrib/binding"
	"bufio"
	"github.com/xiangshouding/martini-middleware/fis"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"
)

const (
	s3 = "-1" //出错
	s0 = "0"  //已安装
	s1 = "1"  //新提交
	s2 = "2"  //安装中
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("public"))

	m.Use(fis.Renderer(fis.Options{
		Directory:  "template",
		Extensions: []string{".tpl"},
	}))

	db, err := sql.Open("mysql", "root@/plg")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	m.Get("/", func(r fis.Render) {
		comps := List(db)
		p := map[string]interface{}{
			"title":      "Submit FIS plugin",
			"components": comps,
		}

		r.HTML(200, "page/index", p)
	})

	m.Post("/new.do", func(res http.ResponseWriter, req *http.Request) {
		plg := req.FormValue("plg")

		if len(strings.TrimSpace(plg)) == 0 {
			http.Redirect(res, req, "/", http.StatusFound)
			// res.Header().Add("location", "/")
			// res.WriteHeader(302)
			return
		}

		comp := strings.Split(strings.TrimSpace(plg), "@")
		var name, version string

		name = comp[0]

		if len(comp) == 2 {
			version = comp[1]
		} else {
			version = "latest"
		}

		New_(db, name, s1, version, "")

		http.Redirect(res, req, "/", http.StatusFound)

	})

	m.Post("/action.do", func(res http.ResponseWriter, req *http.Request) {
		typ := strings.TrimSpace(req.FormValue("type"))
		comp := strings.TrimSpace(req.FormValue("comp"))
		res.Header().Add("content-type", "text/json")
		if typ == "" || comp == "" {
			json_ := json.NewEncoder(res)
			json_.Encode(map[string]string{
				"code": "1",
				"msg":  "require the `name@version` of component",
			})
			return
		}
		code := "1"
		msg := "install fail"
		switch typ {
		case "install":
			code, msg = Install(comp, db, false)
		case "update":
			code, msg = Install(comp, db, true) //安装最新
		}

		json_ := json.NewEncoder(res)
		json_.Encode(map[string]string{
			"code": code,
			"msg":  msg,
		})
	})

	m.Post("/refresh.do", func(r fis.Render) {
		Refresh(db)
		r.JSON(200, map[string]interface{}{"code": "0"})
	})

	fmt.Println(martini.Env)
	m.Run()
}

func List(db *sql.DB) []map[string]interface{} {
	rows, err := db.Query("select * from components order by id desc")

	if err != nil {
		return []map[string]interface{}{}
	}

	new_ := []map[string]interface{}{}

	for rows.Next() {
		var r1 int64
		var r2 string
		var r3 int
		var r4 string
		var r5 string
		var r6 string
		if err := rows.Scan(&r1, &r2, &r3, &r4, &r5, &r6); err != nil {
			log.Fatal(err)
		}

		var row = make(map[string]interface{})
		row["_id"] = r1
		row["name"] = r2
		row["status"] = r3
		row["version"] = r4
		row["user"] = r5
		row["create_date"] = r6

		new_ = append(new_, row)
	}

	return new_
}

func New_(db *sql.DB, name, status, version, user string) (int64, error) {
	var name_ sql.NullString
	err := db.QueryRow("SELECT name FROM components WHERE name=?", name).Scan(&name_)

	if name_.Valid {
		return -1, err
	}

	stmt, err := db.Prepare("INSERT INTO  components VALUES(null, ?, ?, ?, ?, ?)")
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	t_ := time.Now()
	result, err := stmt.Exec(name, status, version, user, t_.Unix())
	if err != nil {
		return -1, err
	}

	last_id, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return last_id, nil
}

func Update_(db *sql.DB, name, status, version, user string) (int64, error) {
	last_id, err := New_(db, name, status, version, user)

	if err == nil {
		return last_id, err
	}

	stmt, err := db.Prepare("UPDATE components SET status=?, version=?, user=? WHERE name=?")
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	result, err := stmt.Exec(version, status, version, user, name)
	if err != nil {
		return -1, err
	}

	last_id, err = result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return last_id, nil
}

func Update_status(db *sql.DB, comp, status string) (bool, int64) {
	result, err := db.Exec("update components set status=? where name=?", status, comp)

	if err != nil {
		return false, -1
	}

	last_id, err := result.LastInsertId()

	if err != nil {
		return false, -1
	}

	return true, last_id
}

func List_local(settings map[string]string) []map[string]string {
	npm_, ok := settings["npm_path"]
	if !ok {
		npm_ = os.Getenv("NODE_PATH")
	}

	dir_arr, err := ioutil.ReadDir(npm_)

	if err != nil {
		return []map[string]string{}
	}

	comps := []map[string]string{}

	for _, dir := range dir_arr {
		if dir.IsDir() {
			package_json := path.Join(npm_, dir.Name(), "package.json")
			content, err := ioutil.ReadFile(package_json)
			if err != nil {
				continue
			}
			decoder := json.NewDecoder(bytes.NewBuffer(content))
			var json_ map[string]interface{}
			err = decoder.Decode(&json_)
			if err != nil {
				log.Println(err)
				continue
			}
			version := json_["version"].(string)
			comp := map[string]string{
				"name":    dir.Name(),
				"version": version,
			}
			comps = append(comps, comp)
		}
	}
	return comps
}

func Refresh(db *sql.DB) bool {
	components := List_local(map[string]string{})

	for _, v := range components {
		Update_(db, v["name"], s0, v["version"], "fis")
	}

	return true
}

func Install(comp string, db *sql.DB, is_update bool) (code, msg string) {
	name := strings.Split(comp, "@")[0]
	Update_status(db, name, "2") // start install
	var cmd *exec.Cmd

	if is_update {
		cmd = exec.Command("npm", "install", "-g", name)
	} else {
		cmd = exec.Command("npm", "install", "-g", comp)
	}

	var stderr = bytes.NewBufferString("")
	var stdout = bytes.NewBufferString("")

	_, err := cmd.StderrPipe()

	if err != nil {
		return "1", err.Error()
	}

	_, err = cmd.StdoutPipe()

	if err != nil {
		return "1", err.Error()
	}

	cmd.Stderr = bufio.NewWriter(stderr)
	cmd.Stdout = bufio.NewWriter(stdout)

	cmd.Start()

	//获取执行错误code真难
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Println(status.ExitStatus())
				return "1", stderr.String()
			}
		}
	}

	Update_status(db, name, "0") //update status

	return "0", stdout.String()
}
