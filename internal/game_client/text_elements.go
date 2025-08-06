package game_client

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

// go:embed space_grotesk.ttf
var fontdata []byte

const (
	TEXT_LOGIN_ANON_BUTTON uint16 = iota
	_textCount
)

const (
	LANG_EN uint32 = iota
	LANG_SP
	LANG_FR
	LANG_DE
	LANG_RU
	LANG_HI
	LANG_JA
	LANG_KO
	LANG_ZH
	LANG_AR
	LANG_IT
	LANG_PT
	LANG_VI
	_langCount
)

var LANG_STR_TABLE = map[string]uint32{
	"en": LANG_EN,
	"sp": LANG_SP,
	"fr": LANG_FR,
	"de": LANG_DE,
	"ru": LANG_RU,
	"hi": LANG_HI,
	"ja": LANG_JA,
	"ko": LANG_KO,
	"zh": LANG_ZH,
	"ar": LANG_AR,
	"it": LANG_IT,
	"pt": LANG_PT,
	"vi": LANG_VI,
}

var AppLang uint32 = LANG_EN
var AppFontFace text.GoTextFace

func GetLangIdx(langString string) uint32 {
	idx, exists := LANG_STR_TABLE[langString]
	if !exists {
		primary, _, _ := strings.Cut(langString, "-")
		idx, exists = LANG_STR_TABLE[primary]
		if !exists {
			idx = LANG_EN
		}
	}
	return idx
}

func SetLangIdx(langString string) {
	AppLang = GetLangIdx(langString)
}

func GetGuiText(textIdx uint16, lang uint8) string {
	return TEXT[textIdx][lang]
}

func GetUiTextSize(size float32, lang uint8, textIdx uint16) (w, h float32) {
	AppFontFace.Size = float64(size)
	t := GetGuiText(textIdx, lang)
	ww, hh := text.Measure(t, &AppFontFace, 1.0)
	return float32(ww), float32(hh)
}

func initFontFace() {
	fontReader := bytes.NewReader(fontdata)
	faceSource, err := text.NewGoTextFaceSource(fontReader)
	if err != nil {
		panic(fmt.Sprintf("could not parse font data"))
	}
	AppFontFace.Source = faceSource
	AppFontFace.Direction = text.DirectionLeftToRight
	AppFontFace.Language = language.English
	AppFontFace.Size = 14.0
}

var TEXT = [_textCount][_langCount]string{
	TEXT_LOGIN_ANON_BUTTON: {
		LANG_EN: "Play Anonymously",
		LANG_SP: "Jugar Anónimamente",
		LANG_FR: "Jouez de Manière Anonyme",
		LANG_DE: "Spielen Sie Anonym",
		LANG_RU: "Играйте анонимно",
		LANG_HI: "गुमनाम रूप से खेलें",
		LANG_JA: "匿名でプレイ",
		LANG_KO: "익명으로 플레이",
		LANG_ZH: "匿名游戏",
		LANG_AR: "العب بشكل مجهول",
		LANG_IT: "Gioca in Modo Anonimo",
		LANG_PT: "Jogue Anonimamente",
		LANG_VI: "Chơi ẩn danh",
	},
}

// var TEXT_USER_DATA = [_textCount]UI_UserData{
// 	TEXT_LOGIN_ANON_BUTTON: UI_UserData{
// 		Kind: USER_DATA_TEXT,
// 		Idx:  uint16(TEXT_LOGIN_ANON_BUTTON),
// 		Color: color.,
// 	},
// }
