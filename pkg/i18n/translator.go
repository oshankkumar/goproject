package i18n

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sudo-suhas/xgo/errors"
	"golang.org/x/text/language"
)

const (
	fileSuffix = ".toml"
)

func MustTranslator(translationsPath string) *Translator {
	t, err := NewTranslator(translationsPath)
	if err != nil {
		panic(fmt.Sprintf("NewTranslator: %v", err))
	}
	return t
}

func NewTranslator(translationsPath string) (*Translator, error) {
	bundler := i18n.NewBundle(language.English)
	bundler.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	files, err := filepath.Glob(path.Join(translationsPath, "*"+fileSuffix))
	if err != nil {
		return nil, errors.E(errors.WithTextf("translations file path: %s", translationsPath), errors.WithErr(err))
	}

	localizers := make(map[string]*i18n.Localizer)
	for _, file := range files {
		if _, err := bundler.LoadMessageFile(file); err != nil {
			return nil, err
		}
		lang := strings.TrimSuffix(filepath.Base(file), fileSuffix)
		localizers[lang] = i18n.NewLocalizer(bundler, mapLanguage(lang))
	}

	return &Translator{bundle: bundler, localizers: localizers}, nil
}

type Translator struct {
	bundle     *i18n.Bundle
	localizers map[string]*i18n.Localizer
}

type TranslationConfig struct {
	Key          string
	TemplateData map[string]interface{}
	PluralCount  interface{}
}

func (t *Translator) Translate(lang string, conf TranslationConfig) string {
	lang = mapLanguage(lang)

	localizer, ok := t.localizers[lang]
	if !ok {
		localizer = i18n.NewLocalizer(t.bundle, lang)
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: conf.Key, TemplateData: conf.TemplateData, PluralCount: conf.PluralCount})
	if err != nil {
		return conf.Key
	}

	return msg
}

func (t *Translator) Title(lang string, conf TranslationConfig) string {
	conf.Key += "_title"
	return t.Translate(lang, conf)
}

func (t *Translator) Message(lang string, conf TranslationConfig) string {
	conf.Key += "_message"
	return t.Translate(lang, conf)
}

func mapLanguage(lang string) string {
	switch lang {
	case "id", "id-ID", "id_ID", "in-ID", "in_ID", "in-id", "in_id":
		return "id-ID"
	case "en-id", "en_id", "en-ID", "en_ID", "en-TH", "en_TH", "en-VN", "en_VN":
		return "en-ID"
	case "th-TH", "th_TH":
		return "th-TH"
	case "vi-VN", "vi_VN", "vi-ID", "vi_ID":
		return "vi-VN"
	default:
		return "en-ID"
	}
}
