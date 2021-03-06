package processAst

import (
	"fmt"
	"github.com/dave/dst"
	"go/token"
	"strconv"
)

func (fld *Fld) Print() *dst.CallExpr  {
	res := &dst.CallExpr{
		Fun: &dst.SelectorExpr{X: &dst.Ident{Name: "t"}, Sel: &dst.Ident{Name: fld.FuncName}},
	}
	args := []dst.Expr{}
	if len(fld.Name) > 0 {
		args = append(args, &dst.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", fld.Name)})
	}
	if len(fld.NameRu) > 0 {
		args = append(args, &dst.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", fld.NameRu)})
	}
	if fld.FuncName == GET_FLD_SELECT_STRING ||  fld.FuncName == GET_FLD_STRING{
		args = append(args, &dst.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", fld.Size)})
	}
	if fld.FuncName == GET_FLD_REF {
		args = append(args, &dst.BasicLit{Kind: token.INT, Value:  fmt.Sprintf("%q", fld.RefTable)})
	}
	// печать RowCol
	if len(fld.RowCol)>0 {
		rows := &dst.CompositeLit{
			Type: &dst.ArrayType{
				Elt: &dst.ArrayType{
					Elt: &dst.Ident{Name: "int"},
				},
			},
			Elts: []dst.Expr{},
		}
		for _, r := range fld.RowCol {
			rows.Elts = append(rows.Elts, &dst.CompositeLit{
				Elts: []dst.Expr {
					&dst.BasicLit{ Kind: token.INT, Value: strconv.Itoa(r[0])},
					&dst.BasicLit{ Kind: token.INT, Value: strconv.Itoa(r[1])},
				},
			},)
		}
		args = append(args, rows)
	}
	if fld.FuncName == GET_FLD_SELECT_STRING || fld.FuncName == GET_FLD_SELECT_MULTIPLE || fld.FuncName == GET_FLD_RADIO_STRING {
		args = append(args, printFldVueOptionsItem(fld))
	}
	if fld.FuncName == GET_FLD_FILES {
		args = append(args, printFldVueFilesParams(fld))
	}
	if fld.FuncName == GET_FLD_IMG || fld.FuncName == GET_FLD_IMG_LIST {
		args = append(args, printFldVueImageParams(fld))
	}
	// печать params
	args = append(args, &dst.BasicLit{Kind: token.STRING, Value:  fmt.Sprintf("%q", fld.ColClass)})
	for i, p := range fld.Params {
		// для GET_FLD_TITLE_COMPUTED первый параметр ставим первым аргументом в функции
		if i == 0 && fld.FuncName == GET_FLD_TITLE_COMPUTED {
			args = append([]dst.Expr{&dst.BasicLit{Kind: token.STRING, Value:  fmt.Sprintf("%q", p)}}, args...)
			continue
		}
		args = append(args, &dst.BasicLit{Kind: token.STRING, Value:  fmt.Sprintf("%q", p)})
	}
	res.Args = args

	// добавляем функции-модификаторы
	for i := len(fld.ModifierList)-1; i >= 0; i-- {
		args := []dst.Expr{}
		for _, v := range fld.ModifierList[i].Args {
			args = append(args, &dst.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v)})
		}
		res = &dst.CallExpr{
			Fun: &dst.SelectorExpr{
				X: res,
				Sel: &dst.Ident{
					Name: fld.ModifierList[i].Name,
				},
			},
			Args: args,
		}
	}

	return res
}

func printFldVueOptionsItem(fld *Fld) *dst.CompositeLit {
	el := &dst.CompositeLit{
		Type: &dst.ArrayType{
			Elt: &dst.SelectorExpr{X: &dst.Ident{Name: "t"}, Sel: &dst.Ident{Name: "FldVueOptionsItem"}},
		},
		Elts: []dst.Expr{},
	}
	for _, v := range fld.FldVueOptionsItem {
		if len(v) == 0 {
			continue
		}
		el1 := &dst.CompositeLit{Elts: []dst.Expr{}}
		for label, value := range v {
			// пропускаем id, которые проставляем во vue
			if label == "id" {
				continue
			}
			el1.Elts = append(el1.Elts, &dst.KeyValueExpr{Key: &dst.Ident{ Name: label}, Value: &dst.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", value)}})
		}
		el.Elts = append(el.Elts, el1)
	}
	return el
}

func printFldVueFilesParams(fld *Fld) *dst.CompositeLit {
	el := &dst.CompositeLit{
		Type: &dst.SelectorExpr {X: &dst.Ident{Name: "t"}, Sel: &dst.Ident{Name: "FldVueFilesParams"}},
		Elts: []dst.Expr{},
	}
	for label, value := range fld.FldVueFilesParams {
		if len(value) == 0 {
			continue
		}
		v := fmt.Sprintf("%q", value)
		if label == "MaxFileSize" {
			v = value
		}
		el.Elts = append(el.Elts, &dst.KeyValueExpr{Key: &dst.Ident{ Name: label}, Value: &dst.BasicLit{Kind: token.STRING, Value: v}})
	}
	return el
}

func printFldVueImageParams(fld *Fld) *dst.CompositeLit {
	el := &dst.CompositeLit{
		Type: &dst.SelectorExpr {X: &dst.Ident{Name: "t"}, Sel: &dst.Ident{Name: "FldVueImgParams"}},
		Elts: []dst.Expr{},
	}
	for label, value := range fld.FldVueImgParams {
		if len(value) == 0 {
			continue
		}
		v := fmt.Sprintf("%q", value)
		if label == "MaxFileSize" || label=="Width" || label=="CanAddUrls" {
			v = value
		}
		el.Elts = append(el.Elts, &dst.KeyValueExpr{Key: &dst.Ident{ Name: label}, Value: &dst.BasicLit{Kind: token.STRING, Value: v}})
	}
	return el
}
