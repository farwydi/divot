package enums

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// 38
var SpecsMap = map[string][]string{
	"WARRIOR": {
		"ARMS",
		"FURY",
		"PROTECTION",
	},
	"PALADIN": {
		"HOLY",
		"PROTECTION",
		"RETRIBUTION",
	},
	"HUNTER": {
		"BEAST_MASTERY",
		"MARKSMANSHIP",
		"SURVIVAL",
	},
	"ROGUE": {
		"ASSASSINATION",
		"OUTLAW",
		"SUBTLETY",
	},
	"PRIEST": {
		"DISCIPLINE",
		"HOLY",
		"SHADOW",
	},
	"DEATH": {
		"KNIGHT_BLOOD",
		"KNIGHT_FROST",
		"KNIGHT_UNHOLY",
	},
	"SHAMAN": {
		"ELEMENTAL",
		"ENHANCEMENT",
		"RESTORATION",
	},
	"MAGE": {
		"ARCANE",
		"FIRE",
		"FROST",
	},
	"WARLOCK": {
		"AFFLICTION",
		"DEMONOLOGY",
		"DESTRUCTION",
	},
	"MONK": {
		"BREWMASTER",
		"MISTWEAVER",
		"WINDWALKER",
	},
	"DRUID": {
		"BALANCE",
		"FERAL",
		"GUARDIAN",
		"RESTORATION",
	},
	"DEMON": {
		"HUNTER_HAVOC",
		"HUNTER_VENGEANCE",
	},
	"EVOKER": {
		"DEVASTATION",
		"PRESERVATION",
	},
}

var Fix = cases.Title(language.English)
