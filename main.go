package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/xuri/excelize/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func isHan(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

func main() {
	viper.SetConfigName("config") //获取配置文件
	viper.AddConfigPath(".")      //添加配置文件所在的路径
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("打开文件失败: %s\n", err)
		os.Exit(1)
	}

	//获取配置文件
	DbHost := viper.GetString("mysql.host")
	DbUsername := viper.GetString("mysql.username")
	DbPassword := viper.GetString("mysql.password")
	DbCharset := viper.GetString("mysql.charset")
	Dbport := viper.GetString("mysql.port")
	DbName := viper.GetString("mysql.dbname")

	ruler := viper.GetStringMap("config.ruler")

	data := []map[string]interface{}{} //初始化数据表

	dbList := []map[string]interface{}{}
	mysqlpath := strings.Join([]string{DbUsername, ":", DbPassword, "@tcp(", DbHost, ":", Dbport, ")/", DbName, "?charset=", DbCharset}, "") //链接配置文件拼接
	db, err := gorm.Open(mysql.Open(mysqlpath), &gorm.Config{})                                                                       //链接数据库
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		os.Exit(1)
	}

	db.Raw("show databases").Scan(&dbList) // 检索库的列表
	fmt.Printf("show dbName: %s\n", dbList)


	file := excelize.NewFile()
	sheet, err := file.NewSheet("敏感信息")
	if err != nil {
		fmt.Printf("创建工作表失败: %s\n", err)
		os.Exit(1)
	}
	file.SetActiveSheet(sheet)
	// 添加标题行
	file.SetCellValue("敏感信息", "A1", "规则名称")
	file.SetCellValue("敏感信息", "B1", "数据库")
	file.SetCellValue("敏感信息", "C1", "表名称")
	file.SetCellValue("敏感信息", "D1", "字段名称")
	file.SetCellValue("敏感信息", "E1", "数据样例")
	row := 2

	for _, dbRow := range dbList {
		for _, dbName := range dbRow {
			dbNameString := dbName.(string)
			if dbNameString == "mysql" || dbNameString == "information_schema"|| dbNameString == "sys"|| dbNameString == "performance_schema"|| dbNameString == "innodb_sys_data"|| dbNameString == "innodb_sys_undo" {
				continue // 跳过这些数据库
			}
			fmt.Println("正在查询库：", dbNameString)

			mysqlpath := strings.Join([]string{DbUsername, ":", DbPassword, "@tcp(", DbHost, ":", Dbport, ")/", dbNameString, "?charset=", DbCharset}, "") //链接配置文件拼接
			db, err := gorm.Open(mysql.Open(mysqlpath), &gorm.Config{})                                                                                           //链接数据库
			if err != nil {
				fmt.Println("连接数据库失败:", err)
				continue
			}

			tableName := []map[string]interface{}{}
			db.Raw("show tables").Scan(&tableName) // 检索表的列表

			for _, v := range tableName { //循环从数据库取出的表map
				for _, s := range v { //循环表map得到键值对
					sString := s.(string) //转换数据库名称为字符串
					fmt.Println("正在查询表：", sString)
					db.Table(sString).Limit(500).Find(&data) //查找数据map
					for _, dataFor := range data {         //循环返回数据map
						for dateListName, dataForOne := range dataFor { //循环单条数据
							for rulerName, rulerFor := range ruler { //循环出整个规则列表
								var dataForOneString string
								switch dataForOne.(type) {
								case string:
									dataForOneString = dataForOne.(string)
								case int:
									dataForOneString = strconv.Itoa(dataForOne.(int))
								case float64:
									dataForOneString = strconv.FormatFloat(dataForOne.(float64), 'f', -1, 64)
								case int32:
									dataForOneString = strconv.Itoa(int(dataForOne.(int32)))
								case int64:
									dataForOneString = strconv.Itoa(int(dataForOne.(int64)))
								}
								matchDigit, _ := regexp.MatchString(rulerFor.(string), dataForOneString)
								if matchDigit {
									file.SetCellValue("敏感信息", fmt.Sprintf("A%d", row), rulerName)
									file.SetCellValue("敏感信息", fmt.Sprintf("B%d", row), dbNameString)
									file.SetCellValue("敏感信息", fmt.Sprintf("C%d", row), sString)
									file.SetCellValue("敏感信息", fmt.Sprintf("D%d", row), dateListName)
									file.SetCellValue("敏感信息", fmt.Sprintf("E%d", row), dataForOne)

									row++
								}
							}
						}
					}
					data = nil
				}
			}
		}
	}

	// Save the Excel file with the current date and time as the filename
	currentTime := time.Now().Format("20060102_150405")
	filename := currentTime + ".xlsx"
	err = file.SaveAs(filename)
	if err != nil {
		fmt.Printf("保存文件失败: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("敏感信息已保存到文件: %s\n", filename)
}
