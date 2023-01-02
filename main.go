package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func main() {

	dsn := "root:super@/go_gorm?parseTime=true"

	dial := mysql.Open(dsn)

	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		// DryRun จะทดสอบแค่คำสั่ง Sql ไม่ทำงานจริง มีค่า default คือ false
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}

	// err = db.Migrator().CreateTable(TestWithModel{})
	// AutoMigrate ไม่ต้องสนใจว่า create แล้วหรือยัง เพราะจะทำการ select เพื่อเช็คดูก่อนว่ามี Table อยู่ไหม
	// ถ้าไม่ใช่ auto ต้องเช็คด้วย db.Migrator().HasTable()
	// err = db.AutoMigrate(Gender{}, TestWithModel{})

	/*

	 */
	if err != nil {
		fmt.Println(err)
	}
	/*
		ถ้าสร้าง Table ซ้ำจะขึ้น error
		Error 1050: Table 'genders' already exists
	*/

	// CreateGender("Female")
	/*
		INSERT INTO `genders` (`name`) VALUES ('Male')
		============================
		{1 Male}
	*/

	// GetGenders()
	/*
		SELECT * FROM `genders` ORDER BY id
		============================
		[{1 Male} {2 Female}]
	*/

	// GetGender(1)
	// GetGenderByName("Female")

	// UpdateGender2(4, "")
	// DeleteGender(4)

	// db.AutoMigrate(TestWithModel{})
	// CreateTestWithModel(0, "Test1")
	// CreateTestWithModel(0, "Test2")
	// CreateTestWithModel(0, "Test3")

	// GetTests()
	// DeleteTest(3)
	// DeleteRealTest(3)
	// GetTests()

	// db.Migrator().CreateTable(Customer{})
	// err = db.AutoMigrate(Gender{}, TestWithModel{}, Customer{})

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// CreateCustomer("Bond", 1)
	// GetCustomers() // [{1 Bond {0 } 1} {2 joke {0 } 2} {3 nam {0 } 2}]

	// GetCustomersWithPreload()
	// [{1 Bond {1 Male} 1} {2 joke {2 Female} 2} {3 nam {2 Female} 2}]

	UpdateGender2(1, " AAA AA ")

}

type SqlLogger struct {
	// ถ้าต้องการ conform to interface นั้นๆ
	// เพียงใส่ Interface ที่ต้องการเป็น attribute ก็สามารถ conform ได้แล้ว
	logger.Interface
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n============================\n", sql)
}

// ถ้าต้องการเปลี่ยนชื่อ Table name ไม่ให้เป็น pluralized ให้ implement interfect Tabler
func (t Test) TableName() string {
	return "MyTEst"
}

type Test struct {
	ID uint

	// ถ้าไม่ต้องการให้ gorm ตั้งชื่อ field ให้
	Name string `gorm:"column:myName;size:20;unique;not null"`
	// CREATE TABLE `genders` (`id` bigint unsigned AUTO_INCREMENT,`myName` longtext,`desc` longtext,PRIMARY KEY (`id`))
	// มักจะใช้เมื่อมี Database อยู่แล้วซึ่งชื่อไม่ตรงตาม condition ของ gorm

	Desc string `gorm:"column:DES_C;type:varchar(50)"`
}

type TestWithModel struct {
	gorm.Model
	Code uint
	Name string
}

func CreateTestWithModel(code uint, name string) {
	test := TestWithModel{Code: code, Name: name}
	db.Create(&test)
}

func GetTests() {
	tests := []TestWithModel{}
	db.Find(&tests)
	for _, t := range tests {
		fmt.Printf("%v|%v\n", t.ID, t.Name)

	}
}

// ** SoftDelete
func DeleteTest(id uint) {
	// เป็น SoftDelete ยังไม่ได้ลบจริงๆ
	db.Delete(&TestWithModel{}, id)
	/*
		UPDATE `test_with_models` SET `deleted_at`='2023-01-01 19:24:04.366' WHERE `test_with_models`.`id`
		= 3 AND `test_with_models`.`deleted_at` IS NULL
	*/
}

// ลบข้อมูลที่เป็น SoftDelete
func DeleteRealTest(id uint) {
	db.Unscoped().Delete(&TestWithModel{}, id)
}

type Gender struct {
	ID   uint   `gorm:"size:20"`
	Name string `gorm:"unique;varchar(50)"`
}

func CreateGender(name string) {
	gender := Gender{Name: name}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)

}

func GetGenders() {
	/*
		db.First()
		// SELECT * FROM users ORDER BY id LIMIT 1;

		db.Take()
		// SELECT * FROM users LIMIT 1;

		db.Last()
		// SELECT * FROM users ORDER BY id DESC LIMIT 1;

		ทั้ง 3 ตัวจะมี gorm.ErrRecordNotFound ไว้ให้ handler

		ถ้าไม่ต้องการ error ใช้ db.Find()
	*/
	genders := []Gender{}
	tx := db.Order("id").Find(&genders)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)

}

func GetGender(id uint) {
	gender := Gender{}
	// First ปกติจะไม่ WHERE ให้
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func GetGenderByName(name string) {
	gender := Gender{}
	// ถ้าต้องการ WHERE ตัวอื่น ใช้ได้ทั้ง db.Find() และตัวอื่น
	// #1
	// tx := db.First(&gender, "name=?", name)
	// #2
	tx := db.Where("name=?", name).Find(&gender)

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func UpdateGender(id uint, name string) {
	/*
		#1
			มี object อยู่โดยไป Query -> แก้ไข -> Save
			db.First(&user)
			user.Name = "jinzhu 2"
			user.Age = 100
			db.Save(&user)
			// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111
			- จะเสีย query เพื่อค้นหาก่อน 1 ที

		#2
			** ถ้าค่าที่ update ตรงกับ zero value จะไม่ update ให้
			ใช้ Model ส่ง Table -> Where -> Update
			// Update with conditions
			db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
			// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;
			- จะไม่เสีย query เพื่อค้นหา

	*/

	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	gender.Name = name
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	GetGender(id)

}
func UpdateGender2(id uint, name string) {
	// ถ้าค่าที่ update ตรงกับ zero value จะไม่ update ให้
	gender := Gender{Name: name}

	// Named Argument
	tx := db.Model(&Gender{}).Where("id=@myID", sql.Named("myID", id)).Updates(gender)
	// tx := db.Model(&Gender{}).Where("id=?", id).Updates(gender)

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	GetGender(id)
}

func DeleteGender(id uint) {
	// วิธีนี้ข้อมูลจะหายไปจริงๆ
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Printf("Delete id %v", id)

	GetGender(id)

}

type Customer struct {
	ID   uint
	Name string
	// ทำ Association โดยการให้ 1 Customer มีเพียง 1 Gender
	// และ 1 Gender จะมีได้หลายคน
	Gender   Gender // อ้างถึง Table Gender
	GenderID uint
	/*
		CREATE TABLE `customers` (`id` bigint unsigned AUTO_INCREMENT,
		`name` longtext,`gender_id` bigint unsigned,PRIMARY KEY (`id`),
		CONSTRAINT `fk_customers_gender` FOREIGN KEY (`gender_id`) REFERENCES `genders`(`id`))
	*/
}

func CreateCustomer(name string, genderID uint) {
	customer := Customer{Name: name, GenderID: genderID}
	tx := db.Create(&customer)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customer)
}

func GetCustomers() {
	customers := []Customer{}
	tx := db.Find(&customers)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customers)

	/*
		[{1 Bond {0 } 1} {2 joke {0 } 2} {3 nam {0 } 2}]
		{0 } << คือ Gender ที่ไม่มี data อยู่
		ถ้าต้องการนำไปใช้จริงต้องใช้ Join และดึง Gender name มาด้วย

		แต่ถ้าเป็น gorm ใช้  Preloading (Eager Loading)
		gorm จะ query genders ขึ้นมาก่อนตาม ID ที่มีอยู่จริง

		SELECT * FROM `genders` WHERE `genders`.`id` IN (1,2)

		และ query customers

		SELECT * FROM `customers`

		ทำให้สามารถเข้าถึงข้อมูลใน gender ได้
		customer.Gender.Name

		แต่ถ้าอยากได้บางสิ่งแต่ไม่รู้จะเขียนยังไงก็ใช้ Raw SQL
		https://gorm.io/docs/sql_builder.html
		Query Raw SQL with Scan

		type Result struct {
		ID   int
		Name string
		Age  int
		}

		var result Result
		db.Raw("SELECT id, name, age FROM users WHERE name = ?", 3).Scan(&result)

		where column อะไรมาก็ได้ และใช้คำสั่ง Scan เพื่อเป็น Type ที่มารับ Query

	*/
}

func GetCustomersWithPreload() {
	customers := []Customer{}
	// ใส่ชื่อ Field ที่ต้องการ Preload
	// tx := db.Preload("Gender").Find(&customers)

	// // ถ้าใช้คำสั่งนี้จะเอา Association มาให้
	tx := db.Preload(clause.Associations).Find(&customers)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	for _, customer := range customers {
		fmt.Printf("%v|%v|%v\n", customer.ID, customer.Name, customer.Gender.Name)
	}
	/*
	   1|Bond|Male
	   2|joke|Female
	   3|nam|Female

	*/

}
