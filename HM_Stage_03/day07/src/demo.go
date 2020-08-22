package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main()  {
	db, err := sql.Open("mysql", "root:123456@(127.0.0.1:3306)/dsy")
	if err != nil {
		fmt.Println("sql.Open: ", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("db.Ping: ", err)
		return
	}

	//test01(db)
	//test02(db)
	test03(db)

}

func test03(db *sql.DB)  {

	var sid, sage int64
	var sname, saddr string
	//rows := db.QueryRow("select * from stu_info")
	////fmt.Println(rows)
	//rows.Scan(&sid, &sname, &sage, &saddr)
	//fmt.Println(sid, sname, sage, saddr)

	rows, err := db.Query("select  * from stu_info")
	if err != nil {
		fmt.Println("db.Query: ", err)
		return
	}
	for rows.Next() {
		rows.Scan(&sid, &sname, &sage, &saddr)
		fmt.Println(sid, sname, sage, saddr)
	}
}

func test02(db *sql.DB)  {
	stu := [2][4] interface{}{{"代书义", 27, "江苏省苏州市昆山市创业路8号"}, {"代义", 26, "江苏省苏州市昆山市创业路8号"}}
	stmt, err := db.Prepare("insert into stu_info values(null, ?, ?, ?)")
	if err != nil {
		fmt.Println("db.Prepare: ", err)
		return
	}
	for _, s := range stu {
		stmt.Exec(s[0], s[1], s[2])
	}
}

func test01(db *sql.DB)  {
	//sql := "SELECT * FROM stu"
	sql := "insert into stu_info values (null, '代书义', 27, '江苏省苏州市昆山市创业路8号')"
	result, err := db.Exec(sql)
	if err != nil {
		fmt.Println("db.Exec: ", err)
		return
	}
	fmt.Printf("%T\n", result)
	fmt.Println(result)

	n, err := result.RowsAffected()
	if err != nil {
		fmt.Println("db.RowsAffected: ", err)
		return
	}
	fmt.Println(n)
}
