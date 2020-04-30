package base

import (
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"
	"resk/infra"
)

var validate *validator.Validate
var translator ut.Translator

func Validate() *validator.Validate {
	Check(validate)
	return validate
}

func Translate() ut.Translator {
	Check(validate)
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	validate = validator.New()
	cn := zh.New()
	uni := ut.New(cn, cn)
	var found bool
	translator, found = uni.GetTranslator("zh")
	if found {
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Error("Not found translator:zh")
	}
}

func ValidateStruct(s interface{}) (err error) {
	err = Validate().Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError);ok {
			log.Error(err)
		}
		errs , ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				log.Error(err.Translate(translator))
			}
		}
	}
	return nil
}