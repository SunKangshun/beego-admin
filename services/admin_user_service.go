package services

import (
	"beego-admin/form_validate"
	"beego-admin/global"
	"beego-admin/models"
	"beego-admin/utils"
	"encoding/base64"
	"errors"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type AdminUserService struct {

}

//根据id获取一条admin_user数据
func (*AdminUserService)GetAdminUserById(id int) *models.AdminUser {
	o := orm.NewOrm()
	adminUser := models.AdminUser{Id:id}
	err := o.Read(&adminUser)
	if err != nil{
		return nil
	}
	return &adminUser
}

//权限检测
func (*AdminUserService)AuthCheck(url string,authExcept map[string]interface{},loginUser *models.AdminUser) bool {
	authUrl := loginUser.GetAuthUrl()
	if utils.KeyInMap(url,authExcept) || utils.KeyInMap(url,authUrl){
		return true
	}else {
		return false
	}
}

//用户登录验证
func (*AdminUserService)CheckLogin(loginForm form_validate.LoginForm,ctx *context.Context) (*models.AdminUser,error) {
	var adminUser models.AdminUser
	o := orm.NewOrm()
	err := o.QueryTable(new(models.AdminUser)).Filter("username",loginForm.Username).Limit(1).One(&adminUser)
	if err != nil{
		return nil,errors.New("用户不存在")
	}

	decodePasswdStr,err := base64.StdEncoding.DecodeString(adminUser.Password)

	if err != nil || !utils.PasswordVerify(loginForm.Password,string(decodePasswdStr)){
		return nil,errors.New("密码错误.")
	}

	if adminUser.Status != 1{
		return nil,errors.New("用户被冻结.")
	}

	ctx.Output.Session(global.LOGIN_USER,adminUser)

	if loginForm.Remember != ""{
		ctx.SetCookie(global.LOGIN_USER_ID,strconv.Itoa(adminUser.Id),7200)
		ctx.SetCookie(global.LOGIN_USER_ID_SIGN,adminUser.GetSignStrByAdminUser(ctx),7200)
	}else{
		ctx.SetCookie(global.LOGIN_USER_ID,ctx.GetCookie(global.LOGIN_USER_ID),-1)
		ctx.SetCookie(global.LOGIN_USER_ID_SIGN,ctx.GetCookie(global.LOGIN_USER_ID_SIGN),-1)
	}

	return &adminUser,nil

}

//获取admin_user 总数
func (*AdminUserService)GetCount() int {
	count,err := orm.NewOrm().QueryTable(new(models.AdminUser)).Count()
	if err != nil{
		return 0
	}
	return int(count)
}

//获取所有adminuser
func (*AdminUserService)GetAllData() []*models.AdminUser {
	var adminUser []*models.AdminUser
	o := orm.NewOrm().QueryTable(new(models.AdminUser))
	_,err := o.All(&adminUser)
	if err != nil{
		return nil
	}else{
		return adminUser
	}
}
