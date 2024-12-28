package sql

import (
	"src/internal/global"
	"src/internal/utils"
	"time"
)

// 用户获取竞赛信息
func SelectCompetitionInfo() []global.CompetitionRequest {
	var competitions []global.CompetitionRequest
	utils.ConnectSql().Table("competitions").Where("is_visible = ?", true).Order("start_at DESC").Find(&competitions)
	return competitions
}

// 管理员获取竞赛信息
func SelectCompetitionInfoAdmin() []global.AdminCompetitionInfoRequest {
	var competitions []global.AdminCompetitionInfoRequest
	utils.ConnectSql().Table("competitions").Find(&competitions)
	return competitions
}

// 管理员获取指定竞赛ID信息
func SelectCompetitionInfoAdminByCid(Cid int) global.AdminCompetitionInfoRequest {
	var competition global.AdminCompetitionInfoRequest
	utils.ConnectSql().Table("competitions").Where("contest_id = ?", Cid).Find(&competition)
	return competition
}

// 管理员删除竞赛
func DeleteCompetition(Cid int) bool {
	result := utils.ConnectSql().Table("competitions").Where("contest_id = ?", Cid).Delete(&global.CompetitionRequest{})
	return result.RowsAffected > 0
}

// 管理员更新/添加竞赛
func UpdateCompetition(req global.AdminCompetitionInfoRequest) error {
	if err := utils.ConnectSql().Table("competitions").Where("contest_id = ?", req.ContestID).Save(&req).Error; err != nil {
		return err
	}
	return nil
}

// 将用户添加至用户-竞赛表
func AddUserCompetition(userId int, competitionId int) error {
	var userInfo global.User
	utils.ConnectSql().Table("users").Where("uid = ?", userId).Find(&userInfo)
	// 当前时间
	nowDateTime := time.Now()
	if err := utils.ConnectSql().Table("user_competitions").Create(
		&global.UserCompetitions{ContestID: competitionId, Uid: userId, Username: userInfo.Username, Join_date: nowDateTime}).Error; err != nil {
		return err
	}
	return nil
}

// 查询指定竞赛参加的所有用户
func SelectUsersCompetition(competitionId int) []global.UserCompetitions {
	var users []global.UserCompetitions
	utils.ConnectSql().Table("user_competitions").Where("contest_id = ?", competitionId).Find(&users)
	return users
}

// 查询用户是否在指定竞赛中
func SelectUserCompetition(userId int, competitionId int) bool {
	var userCompetition global.UserCompetitions
	utils.ConnectSql().Table("user_competitions").Where("uid = ? AND contest_id = ?", userId, competitionId).Find(&userCompetition)
	return userCompetition.Uid != 0
}

// 将用户从用户-竞赛表删除
func DeleteUserCompetition(userId int, competitionId int) error {
	if err := utils.ConnectSql().Table("user_competitions").Where("uid = ? AND contest_id = ?", userId, competitionId).Delete(&global.UserCompetitions{}).Error; err != nil {
		return err
	}
	return nil
}

// 竞赛状态更新
func UpdateCompetitionStatus() error {
	now := time.Now()
	if err := utils.ConnectSql().Table("competitions").Where("start_at <= ? AND end_at >= ?", now, now).Update("enable", true).Error; err != nil {
		return err
	}
	if err := utils.ConnectSql().Table("competitions").Where("end_at < ?", now).
		Update("enable", false).Error; err != nil {
		return err
	}
	return nil
}
