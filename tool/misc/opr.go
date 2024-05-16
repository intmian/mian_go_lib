package misc

import (
	"errors"
	"fmt"
	"os"
)

type OprType int

const (
	OprTypeNull OprType = iota
	OprTypeRename
	OprTypeMove
	OprCreateDir
)

type Opr struct {
	Type      OprType
	SrcAddr   *string
	NameParam *string
	AddrParam *string
}

func (o *Opr) Text() string {
	switch o.Type {
	case OprTypeRename:
		return fmt.Sprintf("重命名 %s 为 %s", *o.SrcAddr, *o.NameParam)
	case OprTypeMove:
		return fmt.Sprintf("移动 %s 到 %s", *o.SrcAddr, *o.AddrParam)
	case OprCreateDir:
		return fmt.Sprintf("创建文件夹 %s", *o.NameParam)
	default:
		return ""
	}
}

func (o *Opr) Do() error {
	switch o.Type {
	case OprTypeRename:
		if o.SrcAddr == nil || o.NameParam == nil {
			return errors.New("OprTypeRename param error")
		}
		f, err := GetFile(*o.SrcAddr)
		if err != nil {
			return errors.Join(errors.New("OprTypeRename GetFile error"), err)
		}
		err = f.Rename(*o.NameParam)
		if err != nil {
			return errors.Join(errors.New("OprTypeRename Rename error"), err)
		}
	case OprTypeMove:
		if o.SrcAddr == nil || o.AddrParam == nil {
			return errors.New("OprTypeMove param error")
		}
		f, err := GetFile(*o.SrcAddr)
		if err != nil {
			return errors.Join(errors.New("OprTypeMove GetFile error"), err)
		}
		err = f.Move(*o.AddrParam)
		if err != nil {
			return errors.Join(errors.New("OprTypeMove Move error"), err)
		}
	case OprCreateDir:
		if o.NameParam == nil {
			return errors.New("OprCreateDir param error")
		}
		err := os.MkdirAll(*o.NameParam, os.ModePerm)
		if err != nil {
			return errors.Join(errors.New("OprCreateDir MakeEmptyFile error"), err)
		}
	default:
		return errors.New("OprType error")
	}
	return nil
}

type Oprs struct {
	oprs []Opr
}

func (o *Oprs) AddRename(srcAddr string, name string) {
	o.oprs = append(o.oprs, Opr{
		Type:      OprTypeRename,
		SrcAddr:   &srcAddr,
		NameParam: &name,
	})
}

func (o *Oprs) AddMove(srcAddr string, addr string) {
	o.oprs = append(o.oprs, Opr{
		Type:      OprTypeMove,
		SrcAddr:   &srcAddr,
		AddrParam: &addr,
	})
}

func (o *Oprs) AddCreateDir(addr string) {
	o.oprs = append(o.oprs, Opr{
		Type:      OprCreateDir,
		NameParam: &addr,
	})
}

func (o *Oprs) ConfirmAndDo() (bool, error) {
	tips := "将进行以下操作：\n"
	tips += o.GetTips()
	tips += "是否任意键继续，或者n键拒绝(n/other)：_"
	var yes string
	err := Input(tips, 1, &yes)
	if err != nil {
		return false, errors.Join(errors.New("Oprs Confirm Input error"), err)
	}
	if yes == "n" || yes == "N" {
		return false, nil
	} else {
		for _, opr := range o.oprs {
			err := opr.Do()
			if err != nil {
				return false, errors.Join(errors.New("Oprs Confirm Do error"), err)
			}
		}
	}
	return true, nil
}

func (o *Oprs) GetTips() string {
	tips := ""
	for _, opr := range o.oprs {
		tips += opr.Text() + "\n"
	}
	return tips
}
