package gorm

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
	"testing"
	"time"
)

// 模型是标准的 struct，由 Go 的基本数据类型、实现了 Scanner 和 Valuer 接口的自定义类型及其指针或别名组成
// 默认情况下，GORM 使用 ID 作为主键，使用结构体名的 蛇形复数 作为表名，字段名的 蛇形 作为列名
// 注：开头要大写（公开） 不然创建不了对应的字段
type User struct {
	// 继承基础模型（id、时间等字段）
	// 软删除：
	// 1. 如果模型包含了一个 gorm.deletedat 字段（gorm.Model 已经包含了该字段)，它将自动获得软删除的能力
	// 2. 当删除操作提交时，会自动变成UPDATE语句去更新deletedAt字段
	gorm.Model
	Name     string
	Age      int
	Birthday time.Time
}

var db *gorm.DB

func init() {
	// 连接数据库
	dsn := "root:admin123@tcp(192.168.1.105:3306)/go_db?charset=utf8mb4&parseTime=true"
	dia, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	db = dia

	// 自动建表
	_ = db.AutoMigrate(&User{})
}

// GORM 允许用户定义的钩子有 BeforeSave, BeforeCreate, AfterSave, AfterCreate 创建记录时将调用这些钩子方法
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("BeforeCreate")
	return
}

// 查询操作会调用
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	fmt.Println("AfterFind")
	return
}

// 更新操作会调用
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// 如果 Role 字段有变更
	if tx.Statement.Changed("Role") {
		return errors.New("role not allowed to change")
	}
	// 如果 Name 或 Role 字段有变更
	if tx.Statement.Changed("Name", "Admin") {
		tx.Statement.SetColumn("Age", 18)
	}
	// 如果任意字段有变更
	if tx.Statement.Changed() {
		tx.Statement.SetColumn("RefreshedAt", time.Now())
	}
	return
}

// 删除操作会调用
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("BeforeDelete")
	return
}

// 插入操作
func TestGormCreate(t *testing.T) {
	var user = User{Name: "tom", Age: 20, Birthday: time.Now()}
	var users = []User{{Name: "tom1"}, {Name: "tom2"}, {Name: "tom3"}}

	// 插入
	tx := db.Create(&user)
	fmt.Printf("影响的条目数: %v\n", tx.RowsAffected)
	fmt.Printf("ID回填: %v\n", user.ID)

	// 创建记录并更新给出的字段
	db.Select("Name", "Age", "CreatedAt").Create(&user)

	// 创建记录并忽略给出的字段
	db.Omit("Name", "Age", "CreatedAt").Create(&user)

	// 批量插入
	db.Create(&users)

	// 分批插入
	db.CreateInBatches(&users, 100)

	// 根据 map 创建记录时 association、回调钩子不会被调用，且主键也不会自动填充
	db.Model(&User{}).Create(map[string]any{"Name": "jinzhu", "Age": 18})
	db.Model(&User{}).Create([]map[string]any{{"Name": "jinzhu_1", "Age": 18}, {"Name": "jinzhu_2", "Age": 20}})

	// 单独设置会话属性 跳过钩子方法
	db.Session(&gorm.Session{SkipHooks: true}).Create(map[string]any{"Name": "jinzhu_", "Age": 18})
}

// 基础查询
func TestGormSearch(t *testing.T) {
	var user User
	var users []User

	// 获取第一条记录（主键升序）SELECT * FROM users ORDER BY id LIMIT 1;
	db.First(&user)
	// 获取一条记录，没有指定排序字段 SELECT * FROM users LIMIT 1;
	db.Take(&user)
	// 获取最后一条记录（主键降序）SELECT * FROM users ORDER BY id DESC LIMIT 1;
	db.Last(&user)

	// map结果映射
	result := map[string]any{}
	db.Model(&User{}).First(&result)
	db.Table("users").First(&result)
	db.Table("users").Take(&result)

	// 结构体结果映射
	type Language struct {
		Code string
		Name string
	}
	db.First(&Language{})

	// 根据主键检索
	// SELECT * FROM users WHERE id = 10;
	db.First(&user, 10)
	// SELECT * FROM users WHERE id = 10;
	db.First(&user, "10")
	// SELECT * FROM users WHERE id IN (1,2,3);
	db.Find(&users, []int{1, 2, 3})

	// 如果主键是字符串（例如像 uuid）
	// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
	db.First(&user, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")

	// 当目标对象有一个主要值时，将使用主键构建条件
	// SELECT * FROM users WHERE id = 10;
	var userResult = User{Model: gorm.Model{ID: 10}}
	db.First(&userResult)
	// SELECT * FROM users WHERE id = 10;
	db.Model(userResult).First(&result)

	// 检索全部对象
	// SELECT * FROM users;
	db.Find(&users)

	// Distinct
	db.Distinct("name", "age").Order("name, age desc").Find(&users)

	// Scan
	//type ScanResult struct {
	//	Name string
	//	Age  int
	//}
	//var scanResult ScanResult
	//db.Table("users").
	//	Select("name", "age").
	//	Where("name = ?", "Antonio").
	//	Scan(&result)
	//db.Raw("SELECT name, age FROM users WHERE name = ?", "Antonio").Scan(&scanResult)
}

// String条件
func TestSearchStringCondition(t *testing.T) {
	var user User
	var users []User

	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;
	db.Where("name = ?", "jinzhu").First(&user)
	// SELECT * FROM users WHERE name <> 'jinzhu';
	db.Where("name <> ?", "jinzhu").Find(&users)
	// SELECT * FROM users WHERE name IN ('jinzhu','jinzhu 2');
	db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
	// SELECT * FROM users WHERE name LIKE '%jin%';
	db.Where("name LIKE ?", "%jin%").Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;
	db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
	// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';
	db.Where("updated_at > ?", "2000-01-01 00:00:00").Find(&users)
	// SELECT * FROM users WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';
	db.Where("created_at BETWEEN ? AND ?", "2000-01-01 00:00:00", "2000-01-08 00:00:00").Find(&users)
}

// Struct & Map 条件
func TestSearchStructAndMapCondition(t *testing.T) {
	var user User
	var users []User
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;
	db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;
	db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
	// SELECT * FROM users WHERE id IN (20, 21, 22);
	db.Where([]int64{20, 21, 22}).Find(&users)

	// 指定结构体查询字段
	// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;
	db.Where(&User{Name: "jinzhu"}, "name", "Age").Find(&users)
	// SELECT * FROM users WHERE age = 0;
	db.Where(&User{Name: "jinzhu"}, "Age").Find(&users)

	// 选择特定字段
	// SELECT name, age FROM users;
	db.Select("name", "age").Find(&users)
	// SELECT name, age FROM users;
	db.Select([]string{"name", "age"}).Find(&users)
	// SELECT COALESCE(age,'42') FROM users;
	_, _ = db.Table("users").Select("COALESCE(age,?)", 42).Rows()
}

// 内联条件
func TestSearchInlineCondition(t *testing.T) {
	var user User
	var users []User
	// SELECT * FROM users WHERE id = 'string_primary_key';
	db.First(&user, "id = ?", "string_primary_key")
	// SELECT * FROM users WHERE name = "jinzhu";
	db.Find(&user, "name = ?", "jinzhu")
	// SELECT * FROM users WHERE name <> "jinzhu" AND age > 20;
	db.Find(&users, "name <> ? AND age > ?", "jinzhu", 20)
	// SELECT * FROM users WHERE age = 20;
	db.Find(&users, User{Age: 20})
	// SELECT * FROM users WHERE age = 20;
	db.Find(&users, map[string]interface{}{"age": 20})
}

// Not条件 Or条件
func TestSearchNotAndOr(t *testing.T) {
	var user User
	var users []User

	// SELECT * FROM users WHERE NOT name = "jinzhu" ORDER BY id LIMIT 1;
	db.Not("name = ?", "jinzhu").First(&user)
	// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");
	db.Not(map[string]interface{}{"name": []string{"jinzhu", "jinzhu 2"}}).Find(&users)
	// SELECT * FROM users WHERE name <> "jinzhu" AND age <> 18 ORDER BY id LIMIT 1;
	db.Not(User{Name: "jinzhu", Age: 18}).First(&user)
	// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;
	db.Not([]int64{1, 2, 3}).First(&user)

	// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
	db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);
	db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2", Age: 18}).Find(&users)
	// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);
	db.Where("name = 'jinzhu'").Or(map[string]interface{}{"name": "jinzhu 2", "age": 18}).Find(&users)
}

// 排序
func TestSearchOrder(t *testing.T) {
	var users []User

	// SELECT * FROM users ORDER BY age desc, name;
	db.Order("age desc, name").Find(&users)
	// SELECT * FROM users ORDER BY age desc, name;
	db.Order("age desc").Order("name").Find(&users)
	// SELECT * FROM users ORDER BY FIELD(id,1,2,3)
	db.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []any{[]int{1, 2, 3}}, WithoutParentheses: true},
	}).Find(&User{})
}

// Limit & Offset
func TestSearchLimitAndOffset(t *testing.T) {
	var users []User
	var users1, users2 []User

	// SELECT * FROM users LIMIT 3;
	db.Limit(3).Find(&users)
	// SELECT * FROM users LIMIT 10; (users1)
	// SELECT * FROM users; (users2)
	db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
	// SELECT * FROM users OFFSET 3;
	db.Offset(3).Find(&users)
	// SELECT * FROM users OFFSET 5 LIMIT 10;
	db.Limit(10).Offset(5).Find(&users)
	// SELECT * FROM users OFFSET 10; (users1)
	// SELECT * FROM users; (users2)
	db.Offset(10).Find(&users1).Offset(-1).Find(&users2)
}

// Group By & Having
func TestSearchGroupAndHaving(t *testing.T) {
	type DataResult struct {
		Date  time.Time
		Total int
	}
	var result DataResult
	var results []DataResult

	// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "group%" GROUP BY `name` LIMIT 1
	db.Model(&User{}).Select("name, sum(age) as total").Where("name LIKE ?", "group%").Group("name").First(&result)

	// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "group"
	db.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("name = ?", "group").Find(&result)

	rows, _ := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()

	rows, _ = db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()

	db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)

	for rows.Next() {

	}
}

// Joins
func TestSearchJoins(t *testing.T) {
	var user []User
	var users []User
	var results []User

	type TempResult struct {
		Name  string
		Email string
	}

	db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&TempResult{})

	// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id
	rows, _ := db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()

	db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)

	// multiple joins with parameter
	db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)

	for rows.Next() {

	}

	// Joins 预加载
	type Company struct {
		Alive bool
	}

	// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id`;
	db.Joins("Company").Find(&users)

	// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id` AND `Company`.`alive` = true;
	db.Joins("Company", db.Where(&Company{Alive: true})).Find(&users)

	// Joins 一个衍生表
	type User struct {
		Id  int
		Age int
	}

	type Order struct {
		UserId     int
		FinishedAt *time.Time
	}

	query := db.Table("order").
		Select("MAX(order.finished_at) as latest").
		Joins("left join user user on order.user_id = user.id").
		Where("user.age > ?", 18).
		Group("order.user_id")

	// SELECT `order`.`user_id`,`order`.`finished_at` FROM `order` join (SELECT MAX(order.finished_at) as latest FROM `order` left join user user on order.user_id = user.id WHERE user.age > 18 GROUP BY `order`.`user_id`) q on order.finished_at = q.latest
	db.Model(&Order{}).
		Joins("join (?) q on order.finished_at = q.latest", query).
		Scan(&results)
}

// 高级查询
func TestGormSeniorSearch(t *testing.T) {
	// 智能选择字段
	type User struct {
		ID     uint
		Name   string
		Age    int
		Gender string
		// 假设后面还有几百个字段...
	}
	type APIUser struct {
		ID   uint
		Name string
	}
	// 查询时会自动选择 `id`, `name` 字段
	// SELECT `id`, `name` FROM `users` LIMIT 10
	db.Model(&User{}).Limit(10).Find(&APIUser{})

	// Locking (FOR UPDATE)
	var users []User
	// SELECT * FROM `users` FOR UPDATE
	db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
	// SELECT * FROM `users` FOR SHARE OF `users`
	db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).Find(&users)
	// SELECT * FROM `users` FOR UPDATE NOWAIT
	db.Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "NOWAIT",
	}).Find(&users)

	// 子查询
	type Pet struct {
		Name string
	}
	db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age = ?", 18).Find(&User{})
	// SELECT * FROM (SELECT `name`,`age` FROM `users`) as u WHERE `age` = 18
	subQuery1 := db.Model(&User{}).Select("name")
	subQuery2 := db.Model(&Pet{}).Select("name")
	db.Table("(?) as u, (?) as p", subQuery1, subQuery2).Find(&User{})
	// SELECT * FROM (SELECT `name` FROM `users`) as u, (SELECT `name` FROM `pets`)

	// Group 条件
	// SELECT * FROM `pizzas` WHERE (pizza = "pepperoni" AND (size = "small" OR size = "medium")) OR (pizza = "hawaiian" AND size = "xlarge")
	//db.Where(
	//	db.Where("pizza = ?", "pepperoni").Where(db.Where("size = ?", "small").Or("size = ?", "medium")),
	//).Or(
	//	db.Where("pizza = ?", "hawaiian").Where("size = ?", "xlarge"),
	//).Find(&Pizza{}).Statement

	// 带多个列的 In
	// SELECT * FROM users WHERE (name, age, role) IN (("jinzhu", 18, "admin"), ("jinzhu 2", 19, "user"));
	db.Where("(name, age, role) IN ?", [][]interface{}{{"jinzhu", 18, "admin"}, {"jinzhu2", 19, "user"}}).Find(&users)

	// 命名参数
	var user User
	// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"
	db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)
	// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu" ORDER BY `users`.`id` LIMIT 1
	db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu"}).First(&user)

	// Find 至 map
	result := map[string]interface{}{}
	db.Model(&User{}).First(&result, "id = ?", 1)
	var results []map[string]interface{}
	db.Table("users").Find(&results)
}

// FirstOrInit
func TestFirstOrInit(t *testing.T) {
	type User struct {
		ID     uint
		Name   string
		Age    int
		Gender string
	}
	var user User

	// 未找到 user，则根据给定的条件初始化一条记录
	db.FirstOrInit(&user, User{Name: "non_existing"})
	// user -> User{Name: "non_existing"}

	// 找到了 `name` = `jinzhu` 的 user
	db.Where(User{Name: "jinzhu"}).FirstOrInit(&user)
	// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	// 找到了 `name` = `jinzhu` 的 user
	db.FirstOrInit(&user, map[string]interface{}{"name": "jinzhu"})
	// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	// 未找到 user，则根据给定的条件以及 Attrs 初始化 user
	db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// user -> User{Name: "non_existing", Age: 20}

	// 未找到 user，则根据给定的条件以及 Attrs 初始化 user
	db.Where(User{Name: "non_existing"}).Attrs("age", 20).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// user -> User{Name: "non_existing", Age: 20}

	// 找到了 `name` = `jinzhu` 的 user，则忽略 Attrs
	db.Where(User{Name: "Jinzhu"}).Attrs(User{Age: 20}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = jinzhu' ORDER BY id LIMIT 1;
	// user -> User{ID: 111, Name: "Jinzhu", Age: 18}

	// 未找到 user，根据条件和 Assign 属性初始化 struct
	db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrInit(&user)
	// user -> User{Name: "non_existing", Age: 20}

	// 找到 `name` = `jinzhu` 的记录，依然会更新 Assign 相关的属性
	db.Where(User{Name: "Jinzhu"}).Assign(User{Age: 20}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = jinzhu' ORDER BY id LIMIT 1;
	// user -> User{ID: 111, Name: "Jinzhu", Age: 20}
}

// FirstOrCreate
func TestFirstOrCreate(t *testing.T) {
	type User struct {
		ID     uint
		Name   string
		Age    int
		Gender string
	}
	var user User

	// 未找到 User，根据给定条件创建一条新纪录
	db.FirstOrCreate(&user, User{Name: "non_existing"})
	// INSERT INTO "users" (name) VALUES ("non_existing");
	// user -> User{ID: 112, Name: "non_existing"}
	// result.RowsAffected // => 1

	// 找到 `name` = `jinzhu` 的 User
	db.Where(User{Name: "jinzhu"}).FirstOrCreate(&user)
	// user -> User{ID: 111, Name: "jinzhu", "Age": 18}
	// result.RowsAffected // => 0

	// 未找到 user，根据条件和 Assign 属性创建记录
	db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
	// user -> User{ID: 112, Name: "non_existing", Age: 20}

	// 找到了 `name` = `jinzhu` 的 user，则忽略 Attrs
	db.Where(User{Name: "jinzhu"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;
	// user -> User{ID: 111, Name: "jinzhu", Age: 18}

	// 未找到 user，根据条件和 Assign 属性创建记录
	db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
	// user -> User{ID: 112, Name: "non_existing", Age: 20}

	// 找到了 `name` = `jinzhu` 的 user，依然会根据 Assign 更新记录
	db.Where(User{Name: "jinzhu"}).Assign(User{Age: 20}).FirstOrCreate(&user)
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;
	// UPDATE users SET age=20 WHERE id = 111;
	// user -> User{ID: 111, Name: "jinzhu", Age: 20}
}

// 优化器、索引提示
func TestHints(t *testing.T) {
	// 优化器提示用于控制查询优化器选择某个查询执行计划，GORM 通过 gorm.io/hints 提供支持
	// SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM `users`
	db.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find(&User{})

	// 索引提示允许传递索引提示到数据库，以防查询计划器出现混乱
	// SELECT * FROM `users` USE INDEX (`idx_user_name`)
	db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
	// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)"
	db.Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find(&User{})
}

// 结果迭代
func TestIterator(t *testing.T) {
	rows, _ := db.Model(&User{}).
		Where("name = ?", "jinzhu").
		Rows()
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var user User
		// ScanRows 方法用于将一行记录扫描至结构体
		_ = db.ScanRows(rows, &user)
		// 业务逻辑...
	}
}

// 用于批量查询并处理记录
func TestFindInBatches(t *testing.T) {
	var results []User
	// 每次批量处理 100 条
	db.Where("processed = ?", false).
		FindInBatches(&results, 100,
			func(tx *gorm.DB, batch int) error {
				for _, result := range results {
					// 批量处理找到的记录
					fmt.Println(result)
				}
				tx.Save(&results)
				// 如果返回错误会终止后续批量操作
				return nil
			})
}

// 用于从数据库查询单个列，并将结果扫描到切片
func TestPluck(t *testing.T) {
	var users []User

	var ages []int64
	db.Model(&users).Pluck("age", &ages)

	var names []string
	db.Model(&User{}).Pluck("name", &names)

	db.Table("deleted_users").Pluck("name", &names)

	// Distinct Pluck
	db.Model(&User{}).Distinct().Pluck("Name", &names)
	// SELECT DISTINCT `name` FROM `users`

	// 超过一列的查询，应该使用 `Scan` 或者 `Find`，例如：
	db.Select("name", "age").Scan(&users)
	db.Select("name", "age").Find(&users)
}

func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "C")
}

func PaidWithCod(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "C")
}

func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status IN (?)", status)
	}
}

// 允许指定常用的查询，可以在调用方法时引用这些查询
func TestScope(t *testing.T) {
	var orders []User

	// 查找所有金额大于 1000 的信用卡订单
	db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)

	// 查找所有金额大于 1000 的货到付款订单
	db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)

	// 查找所有金额大于 1000 且已付款或已发货的订单
	db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
}

// 用于获取匹配的记录数
func TestCount(t *testing.T) {
	var count int64
	db.Model(&User{}).Where("name = ?", "jinzhu").Or("name = ?", "jinzhu 2").Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2'

	db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count)
	// SELECT count(1) FROM users WHERE name = 'jinzhu'; (count)

	db.Table("deleted_users").Count(&count)
	// SELECT count(1) FROM deleted_users;

	// Count with Distinct
	db.Model(&User{}).Distinct("name").Count(&count)
	// SELECT COUNT(DISTINCT(`name`)) FROM `users`

	db.Table("deleted_users").Select("count(distinct(name))").Count(&count)
	// SELECT count(distinct(name)) FROM deleted_users

	// Count with Group
	db.Model(&User{}).Group("name").Count(&count)
}

// 更新操作
func TestUpdate(t *testing.T) {
	var user User
	db.First(&user)
	user.Name = "jinzhu 2"
	user.Age = 100

	// Save 会保存所有的字段，即使字段是零值
	// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111;
	db.Save(&user)

	// 条件更新
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;
	db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;
	db.Model(&user).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
	db.Model(&user).Where("active = ?", true).Update("name", "hello")

	// 更新多列
	// 根据 struct 更新属性，只会更新非零值的字段
	// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;
	db.Model(&user).Updates(User{Name: "hello", Age: 18})
	// 根据 map 更新属性
	// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
	db.Model(&user).Updates(map[string]any{"name": "hello", "age": 18})

	// 更新选定字段
	// UPDATE users SET name='hello' WHERE id=111;
	db.Model(&user).Select("name").Updates(map[string]any{"name": "hello", "age": 18, "active": false})
	// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
	db.Model(&user).Omit("name").Updates(map[string]any{"name": "hello", "age": 18, "active": false})
	// UPDATE users SET name='new_name', age=0 WHERE id=111;
	db.Model(&user).Select("Name", "Age").Updates(User{Name: "new_name", Age: 0})
	// Select all fields (select all fields include zero value fields)
	db.Model(&user).Select("*").Updates(User{Name: "jinzhu", Age: 0})
	// Select all fields but omit Role (select all fields include zero value fields)
	db.Model(&user).Select("*").Omit("Role").Updates(User{Name: "jinzhu", Age: 0})

	// 批量更新 未通过 Model 指定记录的主键，则 GORM 会执行批量更新
	// UPDATE users SET name='hello', age=18 WHERE role = 'admin';
	db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
	// UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);
	db.Table("users").Where("id IN ?", []int{10, 11}).Updates(map[string]any{"name": "hello", "age": 18})

	// 使用 SQL 表达式更新
	// UPDATE "products" SET "price" = price * 2 + 100, "updated_at" = '2013-11-17 21:34:10' WHERE "id" = 3;
	db.Model(&user).Update("price", gorm.Expr("price * ? + ?", 2, 100))
	// UPDATE "products" SET "price" = price * 2 + 100, "updated_at" = '2013-11-17 21:34:10' WHERE "id" = 3;
	db.Model(&user).Updates(map[string]any{"price": gorm.Expr("price * ? + ?", 2, 100)})
	// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3;
	db.Model(&user).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
	// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3 AND quantity > 1;
	db.Model(&user).Where("quantity > 1").UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))

	// 根据子查询进行更新
	// UPDATE "users" SET "company_name" = (SELECT name FROM companies WHERE companies.id = users.company_id);
	db.Model(&user).Update("company_name", db.Model(&User{}).Select("name").Where("companies.id = users.company_id"))
	db.Table("users as u").Where("name = ?", "jinzhu").Update("company_name", db.Table("companies as c").Select("name").Where("c.id = u.company_id"))
	db.Table("users as u").Where("name = ?", "jinzhu").Updates(map[string]interface{}{"company_name": db.Table("companies as c").Select("name").Where("c.id = u.company_id")})

	// 不使用 Hook 和时间追踪
	// UPDATE users SET name='hello' WHERE id = 111;
	db.Model(&user).UpdateColumn("name", "hello")
	// UPDATE users SET name='hello', age=18 WHERE id = 111;
	db.Model(&user).UpdateColumns(User{Name: "hello", Age: 18})
	// UPDATE users SET name='hello', age=0 WHERE id = 111;
	db.Model(&user).Select("name", "age").UpdateColumns(User{Name: "hello", Age: 0})

	// 返回修改行的数据
	// UPDATE `users` SET `salary`=salary * 2,`updated_at`="2021-10-28 17:37:23.19" WHERE role = "admin" RETURNING *
	// users => []User{{ID: 1, Name: "jinzhu", Role: "admin", Salary: 100}, {ID: 2, Name: "jinzhu.2", Role: "admin", Salary: 1000}}
	var users []User
	db.Model(&users).Clauses(clause.Returning{}).Where("role = ?", "admin").Update("salary", gorm.Expr("salary * ?", 2))
	// UPDATE `users` SET `salary`=salary * 2,`updated_at`="2021-10-28 17:37:23.19" WHERE role = "admin" RETURNING `name`, `salary`
	// users => []User{{ID: 0, Name: "jinzhu", Role: "", Salary: 100}, {ID: 0, Name: "jinzhu.2", Role: "", Salary: 1000}}
	db.Model(&users).Clauses(clause.Returning{Columns: []clause.Column{{Name: "name"}, {Name: "salary"}}}).Where("role = ?", "admin").Update("salary", gorm.Expr("salary * ?", 2))
}

// 删除操作
func TestDelete(t *testing.T) {
	var user User
	var users []User

	// 删除一条记录 需要指定主键，否则会触发 批量 Delete
	// DELETE from user where id = 10;
	db.Delete(&user)
	// DELETE from user where id = 10 AND name = "jinzhu";
	db.Where("name = ?", "jinzhu").Delete(&user)

	// 根据主键删除
	// DELETE FROM users WHERE id = 10;
	db.Delete(&User{}, 10)
	// DELETE FROM users WHERE id = 10;
	db.Delete(&User{}, "10")
	// DELETE FROM users WHERE id IN (1,2,3);
	db.Delete(&users, []int{1, 2, 3})

	// 批量删除 如果指定的值不包括主属性，那么 GORM 会执行批量删除，它将删除所有匹配的记录
	// DELETE from emails where email LIKE "%jinzhu%";
	db.Where("email LIKE ?", "%jinzhu%").Delete(&User{})
	// DELETE from emails where email LIKE "%jinzhu%";
	db.Delete(&User{}, "email LIKE ?", "%jinzhu%")

	// 返回删除行的数据 仅适用于支持 Returning 的数据库
	// DELETE FROM `users` WHERE role = "admin" RETURNING *
	// users => []User{{ID: 1, Name: "jinzhu", Role: "admin", Salary: 100}, {ID: 2, Name: "jinzhu.2", Role: "admin", Salary: 1000}}
	db.Clauses(clause.Returning{}).Where("role = ?", "admin").Delete(&users)
	// DELETE FROM `users` WHERE role = "admin" RETURNING `name`, `salary`
	// users => []User{{ID: 0, Name: "jinzhu", Role: "", Salary: 100}, {ID: 0, Name: "jinzhu.2", Role: "", Salary: 1000}}
	db.Clauses(clause.Returning{Columns: []clause.Column{{Name: "name"}, {Name: "salary"}}}).Where("role = ?", "admin").Delete(&users)

	// 查找被软删除的记录
	db.Unscoped().Where("age = 20").Find(&users)

	// 永久删除
	db.Unscoped().Delete(&user)
}

// 原生 SQL 和 SQL 生成器
func TestSQL(t *testing.T) {
	var user User
	var users []User

	// 原生查询 SQL 和 Scan
	db.Raw("SELECT id, name, age FROM users WHERE name = ?", 3).Scan(&user)
	db.Raw("SELECT id, name, age FROM users WHERE name = ?", 3).Scan(&user)
	var age int
	db.Raw("SELECT SUM(age) FROM users WHERE role = ?", "admin").Scan(&age)
	db.Raw("UPDATE users SET name = ? WHERE age = ? RETURNING id, name", "jinzhu", 20).Scan(&users)

	// Exec 原生 SQL
	db.Exec("DROP TABLE users")
	db.Exec("UPDATE orders SET shipped_at = ? WHERE id IN ?", time.Now(), []int64{1, 2, 3})
	db.Exec("UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")

	// 命名参数 支持 sql.NamedArg、map[string]interface{}{} 或 struct 形式的命名参数
	// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"
	db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)
	// SELECT * FROM `users` WHERE name1 = "jinzhu2" OR name2 = "jinzhu2" ORDER BY `users`.`id` LIMIT 1
	db.Where("name1 = @name OR name2 = @name", map[string]any{"name": "jinzhu2"}).First(&user)
	// SELECT * FROM users WHERE name1 = "jinzhu1" OR name2 = "jinzhu2" OR name3 = "jinzhu1"
	db.Raw("SELECT * FROM users WHERE name1 = @name OR name2 = @name2 OR name3 = @name",
		sql.Named("name", "jinzhu1"), sql.Named("name2", "jinzhu2")).Find(&user)
	// UPDATE users SET name1 = "jinzhunew", name2 = "jinzhunew2", name3 = "jinzhunew"
	db.Exec("UPDATE users SET name1 = @name, name2 = @name2, name3 = @name",
		sql.Named("name", "jinzhunew"), sql.Named("name2", "jinzhunew2"))
	// SELECT * FROM users WHERE (name1 = "jinzhu" AND name3 = "jinzhu") AND name2 = "jinzhu2"
	db.Raw("SELECT * FROM users WHERE (name1 = @name AND name3 = @name) AND name2 = @name2",
		map[string]any{"name": "jinzhu", "name2": "jinzhu2"}).Find(&user)
	type NamedArgument struct {
		Name  string
		Name2 string
	}
	// SELECT * FROM users WHERE (name1 = "jinzhu" AND name3 = "jinzhu") AND name2 = "jinzhu2"
	db.Raw("SELECT * FROM users WHERE (name1 = @Name AND name3 = @Name) AND name2 = @Name2",
		NamedArgument{Name: "jinzhu", Name2: "jinzhu2"}).Find(&user)

	// DryRun 模式 在不执行的情况下生成 SQL 及其参数，可以用于准备或测试生成的 SQL
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
	// SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
	fmt.Println(stmt.SQL.String())
	// []interface{}{1}
	fmt.Println(stmt.Vars)

	// ToSQL 返回生成的 SQL 但不执行
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&User{}).
			Where("id = ?", 100).
			Limit(10).
			Order("age desc").
			Find(&[]User{})
	})
	// SELECT * FROM "users" WHERE id = 100 AND "users"."deleted_at" IS NULL ORDER BY age desc LIMIT 10
	fmt.Println(sql)

	// Row
	var name string
	row := db.Table("users").
		Where("name = ?", "jinzhu").
		Select("name", "age").
		Row()
	_ = row.Scan(&name, &age)
	var email string
	row = db.Raw("select name, age, email from users where name = ?", "jinzhu").Row()
	_ = row.Scan(&name, &age, &email)

	// Rows
	rows, _ := db.Model(&User{}).
		Where("name = ?", "jinzhu").
		Select("name, age, email").
		Rows()
	rows, _ = db.Raw("select name, age, email from users where name = ?", "jinzhu").Rows()
	for rows.Next() {
		_ = rows.Scan(&name, &age, &email)
		// 业务逻辑...
	}

	// 将 sql.Rows 扫描至 model
	rows, _ = db.Model(&User{}).
		Where("name = ?", "jinzhu").
		Select("name, age, email").
		Rows()
	for rows.Next() {
		// ScanRows 将一行扫描至 user
		_ = db.ScanRows(rows, &user)
		// 业务逻辑...
	}

	// 连接 在一条 tcp DB 连接中运行多条 SQL (不是事务)
	_ = db.Connection(func(tx *gorm.DB) error {
		tx.Exec("SET my.role = ?", "admin")
		tx.First(&User{})
		return nil
	})

	// GORM 内部使用 SQL builder 生成 SQL
	// 对于每个操作，GORM 都会创建一个 *gorm.Statement 对象，所有的 GORM API 都是在为 statement 添加、修改 子句
	// 最后，GORM 会根据这些子句生成 SQL
}
