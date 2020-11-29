/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/28
   Description :
-------------------------------------------------
*/

package validator

import (
	"errors"
	"reflect"

	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/kataras/iris/v12"

	"github.com/zlyuancn/zapp/core"
)

type Validator struct {
	validateTrans ut.Translator
	validate      *validator.Validate
}

func NewValidator() core.IValidator {
	zh := zhongwen.New()
	vt, _ := ut.New(zh, zh).GetTranslator("zh")

	validate := validator.New()
	_ = zh_translations.RegisterDefaultTranslations(validate, vt)
	return &Validator{
		validateTrans: vt,
		validate:      validate,
	}
}

func (v *Validator) Valid(a interface{}) error {
	err := v.validate.Struct(a)
	return v.translateValidateErr(err)
}

func (v *Validator) ValidField(a interface{}, tag string) error {
	err := v.validate.Var(a, tag)
	return v.translateValidateErr(err)
}
func (v *Validator) Bind(ctx iris.Context, a interface{}) error {
	if err := ctx.ReadBody(a); err != nil {
		return err
	}

	val := reflect.ValueOf(a)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	err := v.Valid(a)
	if err != nil {
		// todo 转为参数错误
		return err
	}
	return nil
}

func (v *Validator) translateValidateErr(err error) error {
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		for _, e := range errs {
			return errors.New(e.Translate(v.validateTrans))
		}
	}
	return nil
}
