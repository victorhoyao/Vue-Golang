package controller

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"BTaskServer/util/Tools"
	"BTaskServer/util/response"
	"BTaskServer/util/validatorTool"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

type IUserWorkController interface {
	GetUserWorkList(c *gin.Context)
	GetMyWorkList(c *gin.Context)
}

type UserWorkController struct {
	DB *gorm.DB
}

func (u UserWorkController) GetMyWorkList(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority == 1 {
		response.AuthError(c, nil, "管理员无结果")
		return
	}

	var getMyWorkQuery model.GetMyWorkQuery
	if !validatorTool.ValidatorQuery[*model.GetMyWorkQuery](c, &getMyWorkQuery) {
		return
	}

	getUserWorkQuery := model.GetUserWorkQuery{
		UserId:    int(usermodel.ID),
		CountDate: getMyWorkQuery.CountDate,
		PageNum:   getMyWorkQuery.PageNum,
		PageSize:  getMyWorkQuery.PageSize,
	}

	// 查询自己某一天的
	if getMyWorkQuery.CountDate != "" {
		getUserDate(u.DB, c, getUserWorkQuery)
		return
	}

	// 查询自己所有日期的
	if getMyWorkQuery.CountDate == "" {
		getUser(u.DB, c, getUserWorkQuery)
		return
	}

}

func (u UserWorkController) GetUserWorkList(c *gin.Context) {
	user, _ := c.Get("user")
	usermodel := user.(model.User)

	if usermodel.Authority != 1 {
		response.AuthError(c, nil, "权限不足")
		return
	}

	var getUserWorkQuery model.GetUserWorkQuery
	if !validatorTool.ValidatorQuery[*model.GetUserWorkQuery](c, &getUserWorkQuery) {
		return
	}

	// 获取全部
	if getUserWorkQuery.UserId == 0 && getUserWorkQuery.CountDate == "" {
		getAll(u.DB, c, getUserWorkQuery)
		return
	}

	// 获取某个用户
	if getUserWorkQuery.UserId != 0 && getUserWorkQuery.CountDate == "" {
		getUser(u.DB, c, getUserWorkQuery)
		return
	}

	// 获取某一天
	if getUserWorkQuery.UserId == 0 && getUserWorkQuery.CountDate != "" {
		getDate(u.DB, c, getUserWorkQuery)
		return
	}

	// 获取某个用户某一天
	if getUserWorkQuery.UserId != 0 && getUserWorkQuery.CountDate != "" {
		getUserDate(u.DB, c, getUserWorkQuery)
		return
	}

}

// 获取全部
func getAll(db *gorm.DB, c *gin.Context, getUserWorkQuery model.GetUserWorkQuery) {
	var userWorkList []model.UserWork
	res := db.Raw("select a.* from userWork as a join (select id from userWork ORDER BY countDate desc,userId limit ?,?) b on a.id = b.id", ((*getUserWorkQuery.PageNum - 1) * *getUserWorkQuery.PageSize), *getUserWorkQuery.PageSize).Scan(&userWorkList)
	if res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	db.Model(model.UserWork{}).Where("1=1").Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *getUserWorkQuery.PageNum,
		"result":  userWorkList,
	}, "获取成功！")
}

// 获取某个用户的
func getUser(db *gorm.DB, c *gin.Context, getUserWorkQuery model.GetUserWorkQuery) {
	var userWorkList []model.UserWork
	res := db.Raw("select a.* from userWork as a join (select id from userWork where userId = ? ORDER BY countDate desc limit ?,?) b on a.id = b.id", getUserWorkQuery.UserId, ((*getUserWorkQuery.PageNum - 1) * *getUserWorkQuery.PageSize), *getUserWorkQuery.PageSize).Scan(&userWorkList)
	if res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	db.Model(model.UserWork{}).Where("userId = ?", getUserWorkQuery.UserId).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *getUserWorkQuery.PageNum,
		"result":  userWorkList,
	}, "获取成功！")
}

// 获取某一天的
func getDate(db *gorm.DB, c *gin.Context, getUserWorkQuery model.GetUserWorkQuery) {
	var userWorkList []model.UserWork
	res := db.Raw("select a.* from userWork as a join (select id from userWork where countDate = ? ORDER BY userId limit ?,?) b on a.id = b.id", getUserWorkQuery.CountDate, ((*getUserWorkQuery.PageNum - 1) * *getUserWorkQuery.PageSize), *getUserWorkQuery.PageSize).Scan(&userWorkList)
	if res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	db.Model(model.UserWork{}).Where("countDate = ?", getUserWorkQuery.CountDate).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *getUserWorkQuery.PageNum,
		"result":  userWorkList,
	}, "获取成功！")
}

// 获取某用户某一天的
func getUserDate(db *gorm.DB, c *gin.Context, getUserWorkQuery model.GetUserWorkQuery) {
	var userWorkList []model.UserWork
	res := db.Raw("select a.* from userWork as a join (select id from userWork where userId = ? and countDate = ? ORDER BY userId limit ?,?) b on a.id = b.id", getUserWorkQuery.UserId, getUserWorkQuery.CountDate, ((*getUserWorkQuery.PageNum - 1) * *getUserWorkQuery.PageSize), *getUserWorkQuery.PageSize).Scan(&userWorkList)
	if res.Error != nil {
		response.ServerBad(c, nil, "获取失败")
		return
	}

	var total int64
	db.Model(model.UserWork{}).Where("userId = ? and countDate = ?", getUserWorkQuery.UserId, getUserWorkQuery.CountDate).Count(&total)

	response.Success(c, gin.H{
		"total":   total,
		"pageNum": *getUserWorkQuery.PageNum,
		"result":  userWorkList,
	}, "获取成功！")
}

// 统计
func countUserWork(db *gorm.DB) {
	go func() {

		codeList := []int{43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53}

		var wg sync.WaitGroup
		for {
			fmt.Println(fmt.Sprintf("统计运行数量%s", Tools.GetDateNowFormat(true)))

			userList, err := getUserList()
			if err != nil {
				fmt.Println("获取用户列表失败")
				time.Sleep(time.Second * 2)
				continue
			}

			today := Tools.GetDateNowFormat(false)

			for _, user := range userList {
				if user.Authority == 1 {
					continue
				}
				fmt.Println(fmt.Sprintf("用户id：%d", user.ID))

				userWork := setUserWorkData(user, today)

				for _, code := range codeList {
					var zCount int64
					var successCount int64
					var failCount int64

					zCountsql := fmt.Sprintf("SELECT count(id) FROM `tasklog` WHERE userKey = '%s' and goodsId = %d and status in (3,4) and sumbmitTime LIKE '%sabc';", user.UserKey, code, today)
					zCountsql1 := strings.ReplaceAll(zCountsql, "abc", "%")

					successCountsql := fmt.Sprintf("SELECT count(id) FROM `tasklog` WHERE userKey = '%s' and goodsId = %d and status = 3 and sumbmitTime LIKE '%sabc';", user.UserKey, code, today)
					successCountsql1 := strings.ReplaceAll(successCountsql, "abc", "%")

					failCountsql := fmt.Sprintf("SELECT count(id) FROM `tasklog` WHERE userKey = '%s' and goodsId = %d and status = 4 and sumbmitTime LIKE '%sabc';", user.UserKey, code, today)
					failCountsql1 := strings.ReplaceAll(failCountsql, "abc", "%")

					db.Raw(zCountsql1).Scan(&zCount)
					db.Raw(successCountsql1).Scan(&successCount)
					db.Raw(failCountsql1).Scan(&failCount)

					switch code {
					case 43:
						userWork.BlbofangCount = int(zCount)
						userWork.BlbofangSucessCount = int(successCount)
						userWork.BlbofangFailCount = int(failCount)
					case 44:
						userWork.BlhuiyuanGouCount = int(zCount)
						userWork.BlhuiyuanGouSucessCount = int(successCount)
						userWork.BlhuiyuanGouFailCount = int(failCount)
					case 45:
						userWork.BliDiggCount = int(zCount)
						userWork.BliDiggSucessCount = int(successCount)
						userWork.BliDiggFailCount = int(failCount)
					case 46:
						userWork.BlfenCount = int(zCount)
						userWork.BlfenSucessCount = int(successCount)
						userWork.BlfenFailCount = int(failCount)
					case 47:
						userWork.BliSLCount = int(zCount)
						userWork.BliSLSucessCount = int(successCount)
						userWork.BliSLFailCount = int(failCount)
					case 48:
						userWork.BlTBCount = int(zCount)
						userWork.BlTBSucessCount = int(successCount)
						userWork.BlTBFailCount = int(failCount)
					case 49:
						userWork.KSDiggCount = int(zCount)
						userWork.KSDiggSucessCount = int(successCount)
						userWork.KSDiggFailCount = int(failCount)
					case 50:
						userWork.KSfenCount = int(zCount)
						userWork.KSfenSucessCount = int(successCount)
						userWork.KSfenFailCount = int(failCount)
					case 51:
						userWork.KSSCCount = int(zCount)
						userWork.KSSCSucessCount = int(successCount)
						userWork.KSSCFailCount = int(failCount)
					case 52:
						userWork.BlgsscCount = int(zCount)
						userWork.BlgsscSucessCount = int(successCount)
						userWork.BlgsscFailCount = int(failCount)
					case 53:
						userWork.BlgsfxCount = int(zCount)
						userWork.BlgsfxSucessCount = int(successCount)
						userWork.BlgsfxFailCount = int(failCount)
					}

				}

				wg.Add(1)
				go saveUserWorkData(&wg, db, userWork)
				wg.Wait()

			}

			fmt.Println(fmt.Sprintf("统计运行数量完毕%s", Tools.GetDateNowFormat(true)))
			time.Sleep(time.Second * 5)
		}
	}()
}

func setUserWorkData(user model.User, today string) model.UserWork {
	newUserWork := model.UserWork{
		CreaTime:          Tools.GetDateNowFormat(true),
		UpdateTime:        Tools.GetDateNowFormat(true),
		CountDate:         today,
		UserId:            user.ID,
		UserKey:           user.UserKey,
		UserName:          user.UserName,
		BliDiggCount:      0,
		BliSLCount:        0,
		BlTBCount:         0,
		BlfenCount:        0,
		BlhuiyuanGouCount: 0,
		BlbofangCount:     0,
		BlgsfxCount:       0,
		BlgsscCount:       0,
		KSDiggCount:       0,
		KSSCCount:         0,
		KSfenCount:        0,

		BliDiggSucessCount:      0,
		BliSLSucessCount:        0,
		BlTBSucessCount:         0, // `gorm:"column:blTBSucessCount;not null;default:0;comment:哔哩投币成功数量" json:"blTBSucessCount" form:"blTBSucessCount"`
		BlfenSucessCount:        0, // `gorm:"column:blfenSucessCount;not null;default:0;comment:哔哩粉成功数量" json:"blfenSucessCount" form:"blfenSucessCount"`
		BlhuiyuanGouSucessCount: 0, // `gorm:"column:blhuiyuanGouSucessCount;not null;default:0;comment:哔哩会员购成功数量" json:"blhuiyuanGouSucessCount" form:"blhuiyuanGouSucessCount"`
		BlbofangSucessCount:     0, // `gorm:"column:blbofangSucessCount;not null;default:0;comment:bili播放成功数量" json:"blbofangSucessCount" form:"blbofangSucessCount"`
		BlgsfxSucessCount:       0, // `gorm:"column:blgsfxSucessCount;not null;default:0;comment:哔哩高速分享成功数量" json:"blgsfxSucessCount" form:"blgsfxSucessCount"`
		BlgsscSucessCount:       0, // `gorm:"column:blgsscSucessCount;not null;default:0;comment:哔哩高速收藏成功数量" json:"blgsscSucessCount" form:"blgsscSucessCount"`
		KSDiggSucessCount:       0, // `gorm:"column:KSDiggSucessCount;not null;default:0;comment:KS点赞成功数量" json:"KSDiggSucessCount" form:"KSDiggSucessCount"`
		KSSCSucessCount:         0, // `gorm:"column:KSSCSucessCount;not null;default:0;comment:KS收藏成功数量" json:"KSSCSucessCount" form:"KSSCSucessCount"`
		KSfenSucessCount:        0, // `gorm:"column:KSfenSucessCount;not null;default:0;comment:KS粉成功数量" json:"KSfenSucessCount" form:"KSfenSucessCount"`

		BliDiggFailCount:      0, // `gorm:"column:bliDiggFailCount;not null;default:0;comment:哔哩赞失败数量" json:"bliDiggFailCount" form:"bliDiggFailCount"`
		BliSLFailCount:        0, // `gorm:"column:bliSLFailCount;not null;default:0;comment:哔哩三连失败数量" json:"bliSLFailCount" form:"bliSLFailCount"`
		BlTBFailCount:         0, // `gorm:"column:blTBFailCount;not null;default:0;comment:哔哩投币失败数量" json:"blTBFailCount" form:"blTBFailCount"`
		BlfenFailCount:        0, // `gorm:"column:blfenFailCount;not null;default:0;comment:哔哩粉失败数量" json:"blfenFailCount" form:"blfenFailCount"`
		BlhuiyuanGouFailCount: 0, // `gorm:"column:blhuiyuanGouFailCount;not null;default:0;comment:哔哩会员购失败数量" json:"blhuiyuanGouFailCount" form:"blhuiyuanGouFailCount"`
		BlbofangFailCount:     0, // `gorm:"column:blbofangFailCount;not null;default:0;comment:bili播放失败数量" json:"blbofangFailCount" form:"blbofangFailCount"`
		BlgsfxFailCount:       0, // `gorm:"column:blgsfxFailCount;not null;default:0;comment:哔哩高速分享失败数量" json:"blgsfxFailCount" form:"blgsfxFailCount"`
		BlgsscFailCount:       0, // `gorm:"column:blgsscFailCount;not null;default:0;comment:哔哩高速收藏失败数量" json:"blgsscFailCount" form:"blgsscFailCount"`
		KSDiggFailCount:       0, // `gorm:"column:KSDiggFailCount;not null;default:0;comment:KS点赞失败数量" json:"KSDiggFailCount" form:"KSDiggFailCount"`
		KSSCFailCount:         0, // `gorm:"column:KSSCFailCount;not null;default:0;comment:KS收藏失败数量" json:"KSSCFailCount" form:"KSSCFailCount"`
		KSfenFailCount:        0, // `gorm:"column:KSfenFailCount;not null;default:0;comment:KS粉失败数量" json:"KSfenFailCount" form:"KSfenFailCount"`
	}

	return newUserWork
}

func saveUserWorkData(wg *sync.WaitGroup, db *gorm.DB, userWork model.UserWork) {
	defer wg.Done()

	var userwork model.UserWork
	res := db.Where("countDate = ? and userId = ?", userWork.CountDate, userWork.UserId).Limit(1).Find(&userwork)
	if res.Error != nil {
		return
	}

	if res.RowsAffected == 0 {
		db.Create(&userWork)
	} else {
		db.Model(&userwork).Updates(map[string]interface{}{
			"userName":   userWork.UserName,
			"updateTime": Tools.GetDateNowFormat(true),
			// 总
			"bliDiggCount":      userWork.BliDiggCount,
			"bliSLCount":        userWork.BliSLCount,
			"blTBCount":         userWork.BlTBCount,
			"blfenCount":        userWork.BlfenCount,
			"blhuiyuanGouCount": userWork.BlhuiyuanGouCount,
			"blbofangCount":     userWork.BlbofangCount,
			"blgsfxCount":       userWork.BlgsfxCount,
			"blgsscCount":       userWork.BlgsscCount,
			"KSDiggCount":       userWork.KSDiggCount,
			"KSSCCount":         userWork.KSSCCount,
			"KSfenCount":        userWork.KSfenCount,
			// 成功
			"bliDiggSucessCount":      userWork.BliDiggSucessCount,
			"bliSLSucessCount":        userWork.BliSLSucessCount,
			"blTBSucessCount":         userWork.BlTBSucessCount,
			"blfenSucessCount":        userWork.BlfenSucessCount,
			"blhuiyuanGouSucessCount": userWork.BlhuiyuanGouSucessCount,
			"blbofangSucessCount":     userWork.BlbofangSucessCount,
			"blgsfxSucessCount":       userWork.BlgsfxSucessCount,
			"blgsscSucessCount":       userWork.BlgsscSucessCount,
			"KSDiggSucessCount":       userWork.KSDiggSucessCount,
			"KSSCSucessCount":         userWork.KSSCSucessCount,
			"KSfenSucessCount":        userWork.KSfenSucessCount,
			// 失败
			"bliDiggFailCount":      userWork.BliDiggFailCount,
			"bliSLFailCount":        userWork.BliSLFailCount,
			"blTBFailCount":         userWork.BlTBFailCount,
			"blfenFailCount":        userWork.BlfenFailCount,
			"blhuiyuanGouFailCount": userWork.BlhuiyuanGouFailCount,
			"blbofangFailCount":     userWork.BlbofangFailCount,
			"blgsfxFailCount":       userWork.BlgsfxFailCount,
			"blgsscFailCount":       userWork.BlgsscFailCount,
			"KSDiggFailCount":       userWork.KSDiggFailCount,
			"KSSCFailCount":         userWork.KSSCFailCount,
			"KSfenFailCount":        userWork.KSfenFailCount,
		})
	}

}

func NewUserWorkController() IUserWorkController {
	db := common.GetDB()
	if err := db.AutoMigrate(&model.UserWork{}); err != nil {
		panic("userWork表迁移失败")
	}
	fmt.Println("userWork表迁移成功")

	countUserWork(db)

	return UserWorkController{DB: db}
}
