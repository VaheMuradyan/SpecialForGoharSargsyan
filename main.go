package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Մոդելներ
type User struct {
	gorm.Model
	Name  string `gorm:"size:255;not null"`
	Email string `gorm:"size:255;uniqueIndex;not null"`
	Age   uint   `gorm:"not null"`
	Posts []Post
}

type Post struct {
	gorm.Model
	Title   string `gorm:"size:255;not null"`
	Content string `gorm:"type:text"`
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
}

func main() {
	dsn := "root:java@tcp(127.0.0.1:3306)/hamalsaran?charset=utf8mb4&parseTime=True&loc=Local"

	// Միանում ենք տվյալների բազային
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Չհաջողվեց միանալ տվյալների բազային")
	}

	// Ավտոմատ միգրացիա (ստեղծում է աղյուսակները)
	err = db.AutoMigrate(&User{}, &Post{})
	if err != nil {
		panic("Միգրացիան ձախողվեց")
	}

	// Ստեղծում ենք նոր օգտատեր
	user := User{
		Name:  "Վահե",
		Email: "vahe@mail.com",
		Age:   25,
	}

	result := db.Create(&user)
	if result.Error != nil {
		panic("Չհաջողվեց ստեղծել օգտատեր")
	}
	fmt.Printf("Ստեղծվել է նոր օգտատեր ID-ով: %d\n", user.ID)

	// Ստեղծում ենք պոստ այդ օգտատիրոջ համար
	post := Post{
		Title:   "Առաջին պոստ",
		Content: "Սա իմ 2-րդ պոստն է",
		UserID:  user.ID,
	}

	db.Create(&post)

	// Տարբեր հարցումների օրինակներ

	// Գտնել օգտատեր ըստ ID-ի
	var foundUser User
	db.First(&foundUser, user.ID)
	fmt.Printf("Գտնված օգտատեր: %v\n", foundUser.Name)

	// Գտնել բոլոր օգտատերերին
	var users []User
	db.Find(&users)

	// Գտնել ըստ պայմանի
	var youngUsers []User
	db.Where("age < ?", 30).Find(&youngUsers)

	// Թարմացնել օգտատիրոջ տվյալները
	db.Model(&user).Updates(User{
		Name: "Նոր Անուն",
		Age:  26,
	})

	// Գտնել օգտատիրոջ բոլոր պոստերը
	var userPosts []Post
	db.Where("user_id = ?", user.ID).Find(&userPosts)

	// Ջնջել պոստը (soft delete)
	db.Delete(&post)

	// Ջնջել օգտատիրոջը (soft delete)
	db.Delete(&user)
}

func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	return result.Error
}

func GetUserByID(db *gorm.DB, id uint) (User, error) {
	var user User
	result := db.First(&user, id)
	return user, result.Error
}

func UpdateUser(db *gorm.DB, user *User) error {
	result := db.Save(user)
	return result.Error
}

func DeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&User{}, id)
	return result.Error
}

// Բարդ հարցումների օրինակներ
func GetUsersWithPosts(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Preload("Posts").Find(&users)
	return users, result.Error
}

func GetUsersByAgeRange(db *gorm.DB, minAge, maxAge uint) ([]User, error) {
	var users []User
	result := db.Where("age BETWEEN ? AND ?", minAge, maxAge).Find(&users)
	return users, result.Error
}
