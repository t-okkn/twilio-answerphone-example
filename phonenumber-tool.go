package main

import (
	"errors"
	"strconv"
	"strings"
)

type PhoneNumber string

// summary => E.164形式の電話番号か確認します
// param::pn => PhoneNumber
// return::bool => E.164形式かどうか
/////////////////////////////////////////
func (pn PhoneNumber) IsE164() bool {
	p := pn.removeHyphen()
	
	if strings.HasPrefix(p, "+") {
		if _, err := strconv.Atoi(p[1:]); err != nil {
			return false
		}
		
		if len([]rune(p)) == 12 || len([]rune(p)) == 13 {
			return true
		}
	}
	
	return false
}

// summary => 0AB-J形式の電話番号か確認します
// param::pn => PhoneNumber
// return::bool => 0AB-J形式かどうか
/////////////////////////////////////////
func (pn PhoneNumber) Is0ABJ() bool {
	p := pn.removeHyphen()
	
	if _, err := strconv.Atoi(p); err != nil {
		return false
	}
	
	if len([]rune(p)) == 10 || len([]rune(p)) == 11 {
		return true
	}
	
	return false
}
	

// summary => E.164形式の電話番号に変換します
// param::pn => PhoneNumber
// return::string => E.164形式の電話番号
// return::error => エラー
/////////////////////////////////////////
func (pn PhoneNumber) ToE164() (string, error) {
	p := pn.removeHyphen()
	
	if pn.IsE164() {
		return p, nil
	}
	
	if pn.Is0ABJ() {
		return "+81" + p[1:], nil
	}
	
	e := errors.New("電話番号ではない可能性があります")
	return "", e
}

// summary => 0AB-J形式の電話番号に変換します
// param::pn => PhoneNumber
// return::string => 0AB-J形式の電話番号
// return::error => エラー
/////////////////////////////////////////
func (pn PhoneNumber) To0ABJ() (string, error) {
	p := pn.removeHyphen()
	
	if pn.Is0ABJ() {
		return p, nil
	}
	
	if pn.IsE164() {
		return strings.Replace(p, "+81", "0", 1), nil
	}
	
	e := errors.New("電話番号ではない可能性があります")
	return "", e
}

// summary => 電話番号のハイフンを除去します
// param::pn => PhoneNumber
// return::string => ハイフンを除去した電話番号
/////////////////////////////////////////
func (pn PhoneNumber) removeHyphen() string {
	p := string(pn)
	
	if strings.Contains(p, "-") {
		p = strings.Join(strings.Split(p, "-"), "")
	}
	
	return p
}

