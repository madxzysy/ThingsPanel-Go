package controllers

import (
	gvalid "ThingsPanel-Go/initialize/validate"
	"ThingsPanel-Go/models"
	"ThingsPanel-Go/services"
	"ThingsPanel-Go/utils"
	valid "ThingsPanel-Go/validate"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	context2 "github.com/beego/beego/v2/server/web/context"
)

type TpOtaController struct {
	beego.Controller
}

// 列表
func (TpOtaController *TpOtaController) List() {
	PaginationValidate := valid.TpOtaPaginationValidate{}
	err := json.Unmarshal(TpOtaController.Ctx.Input.RequestBody, &PaginationValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(PaginationValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(PaginationValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			utils.SuccessWithMessage(1000, message, (*context2.Context)(TpOtaController.Ctx))
			break
		}
		return
	}
	var TpOtaService services.TpOtaService
	isSuccess, d, t := TpOtaService.GetTpOtaList(PaginationValidate)

	if !isSuccess {
		utils.SuccessWithMessage(1000, "查询失败", (*context2.Context)(TpOtaController.Ctx))
		return
	}
	dd := valid.RspTpOtaPaginationValidate{
		CurrentPage: PaginationValidate.CurrentPage,
		Data:        d,
		Total:       t,
		PerPage:     PaginationValidate.PerPage,
	}
	utils.SuccessWithDetailed(200, "success", dd, map[string]string{}, (*context2.Context)(TpOtaController.Ctx))

}

// 新增
func (TpOtaController *TpOtaController) Add() {
	AddTpOtaValidate := valid.AddTpOtaValidate{}
	err := json.Unmarshal(TpOtaController.Ctx.Input.RequestBody, &AddTpOtaValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(AddTpOtaValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(AddTpOtaValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			utils.SuccessWithMessage(1000, message, (*context2.Context)(TpOtaController.Ctx))
			break
		}
		return
	}

	// 判断文件是否存在|| !utils.FileExist(AddTpOtaValidate.PackageUrl)
	path := "." + utils.GetUrlPath(AddTpOtaValidate.PackageUrl)
	if !utils.FileExist(path) {
		utils.SuccessWithMessage(400, "升级包不存在", (*context2.Context)(TpOtaController.Ctx))
		return
	}
	if err := utils.CheckPathFilename(path); err != nil || AddTpOtaValidate.PackageUrl == "" {
		utils.SuccessWithMessage(400, "升级包路径不合法或升级包路径是空", (*context2.Context)(TpOtaController.Ctx))
		return
	}
	//文件sign计算
	packagesign, sign_err := utils.FileSign(path, AddTpOtaValidate.SignatureAlgorithm)
	if sign_err != nil {
		utils.SuccessWithMessage(400, "文件签名计算失败", (*context2.Context)(TpOtaController.Ctx))
		return
	}

	//文件大小检查
	packageLength, pl_err := utils.GetFileSize(path)
	if pl_err != nil {
		utils.SuccessWithMessage(400, "文件大小计算失败", (*context2.Context)(TpOtaController.Ctx))
		return
	}
	if packageLength > 1024*1024*1024*1024 {
		utils.SuccessWithMessage(400, "文件大小超出1G", (*context2.Context)(TpOtaController.Ctx))
		return
	}
	var TpOtaService services.TpOtaService
	id := utils.GetUuid()
	TpOta := models.TpOta{
		Id:                 id,
		PackageName:        AddTpOtaValidate.PackageName,
		PackageVersion:     AddTpOtaValidate.PackageVersion,
		PackageModule:      AddTpOtaValidate.PackageModule,
		ProductId:          AddTpOtaValidate.ProductId,
		SignatureAlgorithm: AddTpOtaValidate.SignatureAlgorithm,
		PackageUrl:         AddTpOtaValidate.PackageUrl,
		Description:        AddTpOtaValidate.Description,
		AdditionalInfo:     AddTpOtaValidate.AdditionalInfo,
		CreatedAt:          time.Now().Unix(),
		Sign:               packagesign,
		FileSize:           utils.FormatFileSize(packageLength),
	}
	d, rsp_err := TpOtaService.AddTpOta(TpOta)
	if rsp_err == nil {
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(TpOtaController.Ctx))
	} else {
		var err string
		isTrue := strings.Contains(rsp_err.Error(), "23505")
		if isTrue {
			err = "有值不能重复！"
		} else {
			err = rsp_err.Error()
		}
		utils.SuccessWithMessage(400, err, (*context2.Context)(TpOtaController.Ctx))
	}
}

// 删除
func (TpOtaController *TpOtaController) Delete() {
	TpOtaIdValidate := valid.TpOtaIdValidate{}
	err := json.Unmarshal(TpOtaController.Ctx.Input.RequestBody, &TpOtaIdValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(TpOtaIdValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(TpOtaIdValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			utils.SuccessWithMessage(1000, message, (*context2.Context)(TpOtaController.Ctx))
			break
		}
		return
	}
	if TpOtaIdValidate.Id == "" {
		utils.SuccessWithMessage(1000, "id不能为空", (*context2.Context)(TpOtaController.Ctx))
		return
	}
	var TpOtaService services.TpOtaService
	TpOta := models.TpOta{
		Id: TpOtaIdValidate.Id,
	}
	rsp_err := TpOtaService.DeleteTpOta(TpOta)
	if rsp_err == nil {
		utils.SuccessWithMessage(200, "success", (*context2.Context)(TpOtaController.Ctx))
	} else {
		utils.SuccessWithMessage(400, rsp_err.Error(), (*context2.Context)(TpOtaController.Ctx))
	}
}
