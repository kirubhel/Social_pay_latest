package filter

import (
	"fmt"
	"strings"
)


type FilterItem interface {
	Build() (string ,[]interface{} ,error )
}



func (f Field) Build()(string,[]interface{},error) {


	op:=strings.ToUpper(f.Operator)


	switch op {
	case "IS NULL","IS NOT NULL" :
		return fmt.Sprintf("%s %s",f.Name,f.Operator) ,nil,nil

	case "BETWEEN":
		values, ok := f.Value.([]interface{})
		if !ok || len(values) != 2 {
			return "", nil, fmt.Errorf("BETWEEN operator requires exactly two values")
		}
		return fmt.Sprintf("%s BETWEEN ? AND ?", f.Name), values, nil

	case "IN","NOT IN":
		values, ok := f.Value.([]interface{})
		if !ok || len(values) == 0 {
			return "", nil, fmt.Errorf("%s operator requires a slice of values", op)
		}
		placeholders := strings.Repeat("?,", len(values))
		placeholders = strings.TrimRight(placeholders, ",")
		return fmt.Sprintf("%s %s (%s)", f.Name, op, placeholders), values, nil

	default:
		return fmt.Sprintf("%s %s ?", f.Name, f.Operator), []interface{}{f.Value}, nil

	}
}

// building search query clause
// func (s *Search) Build() (string, []interface{}, error) {
// 	if s == nil || s.Term == "" || len(s.Fields) == 0 {
// 		return "", nil, nil
// 	}

// 	var clauses []string
// 	var args []interface{}

// 	for _, field := range s.Fields {
// 		clauses = append(clauses, fmt.Sprintf("(%s ILIKE ?)", field))
// 		args = append(args, "%"+s.Term+"%")
// 	}

// 	joined := strings.Join(clauses, " OR ")
// 	return fmt.Sprintf("(%s)", joined), args, nil
// }


func (s *Search) Build() (string, []interface{}, error) {
	if s == nil || len(s.Queries) == 0 {
		return "", nil, nil
	}

	var clauses []string
	var args []interface{}

	for _, q := range s.Queries {
		if q.Term == "" || q.Field == "" {
			continue
		}
		clauses = append(clauses, fmt.Sprintf("(%s ILIKE ?)", q.Field))
		args = append(args, "%"+q.Term+"%")
	}

	if len(clauses) == 0 {
		return "", nil, nil
	}

	joined := strings.Join(clauses, " OR ")
	return fmt.Sprintf("(%s)", joined), args, nil
}


// build multiple filter group
func (g FilterGroup) Build()(string,[]interface{},error) {

	if len(g.Fields) == 0 {

		return "",nil,nil 
	}

	var clauses [] string 
	var args [] interface{}

	// building the nested query	
	for _,item:=range g.Fields {
		
		sqlPart,partArgs,err:=item.Build()

		if err != nil {

			return "",nil,err
		}

		if sqlPart != "" {
			clauses = append(clauses, fmt.Sprintf("(%s)", sqlPart))
			args = append(args, partArgs...)
		}

	}

	joined:=strings.Join(clauses,fmt.Sprintf(" %s ",strings.ToUpper(g.Linker)))
	return joined,args,nil
}





func (f Filter) Build()(string,[]interface{},error) {

	whereClause ,args,err:=f.Group.Build()

	if err != nil {
		return "",nil,err
	}

	// search clause 

	searchClause,searchArgs,err:=f.Search.Build()

	if err != nil {

		return "",nil,err
	}

	if searchClause!=""{
		if whereClause != ""{
			whereClause = fmt.Sprintf("%s AND %s", whereClause, searchClause)
		} else {
			whereClause=searchClause
		}

		args=append(args, searchArgs...)
	}


	sortClause := ""
	if len(f.Sort) > 0 {
		var sortParts []string
		for _, s := range f.Sort {
			if s.Field == "" {
				continue
			}
			op := strings.ToUpper(s.Operator)
			if op != "ASC" && op != "DESC" {
				op = "ASC"
			}
			sortParts = append(sortParts, fmt.Sprintf("%s %s", s.Field, op))
		}
		if len(sortParts) > 0 {
			sortClause = " "+"ORDER BY " + strings.Join(sortParts, ", ")
		}
	}

	paginationClause:=""

	if f.Pagination.GetLimit() >0 {

		paginationClause += fmt.Sprintf(" LIMIT %d", f.Pagination.GetLimit())

	}

	if f.Pagination.GetOffset() >0 {
		paginationClause += fmt.Sprintf(" OFFSET %d", f.Pagination.GetOffset())
	}


	if whereClause == ""{
		// no filter
		return "WHERE 1=1",nil, nil
	}

	clause:="WHERE "+ whereClause + sortClause  + paginationClause
	finalSQL, err := convertPlaceholders(clause, args)

	if err != nil {
		return "", nil, err
	}

	return finalSQL + ";", args, nil
}



func convertPlaceholders(sql string, args []interface{}) (string, error) {
	var result strings.Builder
	index := 1

	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' {
			result.WriteString(fmt.Sprintf("$%d", index))
			index++
		} else {
			result.WriteByte(sql[i])
		}
	}

	if index-1 != len(args) {
		return "", fmt.Errorf(" mismatched placeholders and arguments: got %d placeholders, %d args", index-1, len(args))
	}

	return result.String(), nil

}