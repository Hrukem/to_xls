func GetFileTemplateCompetitiveListHandler(w http.ResponseWriter, r *http.Request) {
	var res ResultInfo
	defer func() {
		if rec := recover(); rec != nil {
			report.Error(r, fmt.Errorf("recover panic: %v\n%#v", rec, report.FilesWithLineNum()))
			res.Message = &UserError
			service.ReturnErrorJSON(w, &res, http.StatusInternalServerError)
		}
	}()
	vars := mux.Vars(r)
	res.User = *Auth(r)
	idCompetitive, err := strconv.ParseInt(vars[`id`], 10, 32)
	if err != nil {
		message := `Неверный параметр id.`
		res.Message = &message
		service.ReturnJSON(w, &res)
		return
	}
	buf, err := res.GetFileTemplateCompetitiveList(r, uint(idCompetitive))
	if err != nil {
		res.Done = false
		m := err.Error()
		res.Message = &m
		res.Items = nil
		service.ReturnJSON(w, &res)
		return
	}
	if buf == nil {
		res.Done = true
		m := `Список пуст.`
		res.Message = &m
		res.Items = nil
		service.ReturnJSON(w, &res)
		return
	}

	filename := `attachment; filename="` + time.Now().Format(`2006-01-02 15:04:05`) + `.xlsx"`
	w.Header().Set("Content-Disposition", filename)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Write(buf)

	service.ReturnJSON(w, &res)
}

unc (result *ResultInfo) GetFileTemplateCompetitiveList(r *http.Request, idCompetitiveGroup uint) ([]byte, error) {
	result.Done = false
	db := config.Db.ConnSQLx
	ctx := r.Context()

	type S struct {
		Id                   string
		UidCompetitiveGroup  string  `db:"uid_competitive_group"`
		NameCompetitiveGroup *string `db:"name_competitive_group"`
		AdmissionVolume      *string `db:"admission_volume"`
		CountFirstStep       *string `db:"count_first_step"`
		CountSecondStep      *string `db:"count_second_step"`
		Changed              *string `db:"changed"`
		Uid                  *string
		UidEpgu              *string `db:"uid_epgu"`
		Guid                 *string
		Rating               *string
		WithoutTests         *string `db:"without_tests"`
		ReasonWithoutTests   *string `db:"reason_without_tests"`
		EntranceTest1        *string `db:"test_name1"`
		Result1              *string
		EntranceTest2        *string `db:"test_name2"`
		Result2              *string
		EntranceTest3        *string `db:"test_name3"`
		Result3              *string
		EntranceTest4        *string `db:"test_name4"`
		Result4              *string
		EntranceTest5        *string `db:"test_name5"`
		Result5              *string
		AchievementsMark     *string `db:"mark"`
		EntranceTestMark     *string `db:"sum_point_entr_test"`
		SumMark              *string `db:"sum_point_all"`
		Benefit              *string
		ReasonBenefit        *string `db:"reason_benefit"`
		Agree                *string
		Original             *string
		Addition             *string
		IdStageEnrollment    *string `db:"id_stage_enrollment,omitempty"`
	}

	query := digest.QueryGetTemplateCompetitiveList

	params := []any{idCompetitiveGroup}

	rows, err := db.QueryxContext(ctx, query, params...)
	if err != nil {
		if err = report.QueryError(ctx, r, err, nil, query, params...); err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	var items []S
	for rows.Next() {
		item := S{}
		if err := rows.StructScan(&item); err != nil {
			report.Error(r, err)
			return nil, report.ErrUser
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		report.Error(r, err)
		return nil, report.ErrUser
	}

	if len(items) == 0 {
		return nil, nil
	}

	// result.Done = true
	// result.Items = listTemplates
	f := excelize.NewFile()
	sheet := "Data"
	var index int
	if f.SheetCount == 0 {
		index = f.NewSheet(sheet)
	} else {
		f.SetSheetName(f.GetSheetName(1), sheet)
		index = 1
	}
	styleTitle, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#CCFFFF"],"pattern":1}}`)

	f.SetColWidth(sheet, "A", "AD1", 20)
	// f.SetColWidth(sheet, "B", "B", 50)
	// f.SetColWidth(sheet, "D", "D", 50)
	// f.SetColWidth(sheet, "F", "H", 40)
	// f.SetColWidth(sheet, "S", "S", 50)

	// Set title
	f.SetCellStyle(sheet, `A1`, `AD1`, styleTitle)

	f.SetCellValue(sheet, `A1`, `UidCompetitiveGroup`)
	f.SetCellValue(sheet, `B1`, `NameCompetitiveGroup`)
	f.SetCellValue(sheet, `C1`, `AdmissionVolume`)
	f.SetCellValue(sheet, `D1`, `CountFirstStep`)
	f.SetCellValue(sheet, `E1`, `CountSecondStep`)
	f.SetCellValue(sheet, `F1`, `Changed`)
	f.SetCellValue(sheet, `G1`, `UidEpgu`)
	f.SetCellValue(sheet, `H1`, `Uid`)
	f.SetCellValue(sheet, `I1`, `Rating`)
	f.SetCellValue(sheet, `J1`, `WithoutTests`)
	f.SetCellValue(sheet, `K1`, `ReasonWithoutTests`)
	f.SetCellValue(sheet, `L1`, `EntranceTest1`)
	f.SetCellValue(sheet, `M1`, `Result1`)
	f.SetCellValue(sheet, `N1`, `EntranceTest2`)
	f.SetCellValue(sheet, `O1`, `Result2`)
	f.SetCellValue(sheet, `P1`, `EntranceTest3`)
	f.SetCellValue(sheet, `Q1`, `Result3`)
	f.SetCellValue(sheet, `R1`, `EntranceTest4`)
	f.SetCellValue(sheet, `S1`, `Result4`)
	f.SetCellValue(sheet, `T1`, `EntranceTest5`)
	f.SetCellValue(sheet, `U1`, `Result5`)
	f.SetCellValue(sheet, `V1`, `AchievementsMark`)
	f.SetCellValue(sheet, `W1`, `EntranceTestMark`)
	f.SetCellValue(sheet, `X1`, `SumMark`)
	f.SetCellValue(sheet, `Y1`, `Benefit`)
	f.SetCellValue(sheet, `Z1`, `ReasonBenefit`)
	f.SetCellValue(sheet, `AA1`, `Agreed`)
	f.SetCellValue(sheet, `AB1`, `Original`)
	f.SetCellValue(sheet, `AC1`, `Addition`)
	f.SetCellValue(sheet, `AD1`, `IdStageEnrollment`)

	fn := func(v *string) string {
		if v != nil {
			return *v
		} else {
			return ""
		}
	}
	var last int
	for i, v := range items {
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", i+2), v.UidCompetitiveGroup)
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", i+2), fn(v.NameCompetitiveGroup))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", i+2), fn(v.AdmissionVolume))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", i+2), fn(v.CountFirstStep))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", i+2), fn(v.CountSecondStep))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "F", i+2), fn(v.Changed))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "G", i+2), fn(v.UidEpgu))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "H", i+2), fn(v.Uid))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "I", i+2), fn(v.Rating))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "J", i+2), fn(v.WithoutTests))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "K", i+2), fn(v.ReasonWithoutTests))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "L", i+2), fn(v.EntranceTest1))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "M", i+2), fn(v.Result1))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "N", i+2), fn(v.EntranceTest2))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "O", i+2), fn(v.Result2))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "P", i+2), fn(v.EntranceTest3))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "Q", i+2), fn(v.Result3))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "R", i+2), fn(v.EntranceTest4))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "S", i+2), fn(v.Result4))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "T", i+2), fn(v.EntranceTest5))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "U", i+2), fn(v.Result5))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "V", i+2), fn(v.AchievementsMark))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "W", i+2), fn(v.EntranceTestMark))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "X", i+2), fn(v.SumMark))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "Y", i+2), fn(v.Benefit))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "Z", i+2), fn(v.ReasonBenefit))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "AA", i+2), fn(v.Agree))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "AB", i+2), fn(v.Original))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "AC", i+2), fn(v.Addition))
		f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "AD", i+2), fn(v.IdStageEnrollment))
		last = i
	}
	last++
	// Последняя строка с описаниями полей для конвертера:
	styleFooter, _ := f.NewStyle(`{"fill":{"type":"pattern","color":["#b0ddaa"],"pattern":1}}`)
	f.SetCellStyle(sheet, `A`+strconv.Itoa(last+2), `AD`+strconv.Itoa(last+2), styleFooter)

	var i = 0
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "A", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "B", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "C", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "D", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "E", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "F", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "G", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "H", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "I", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "J", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "K", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "L", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "M", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "N", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "O", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "P", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "Q", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "R", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "S", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "T", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "U", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "V", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "W", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "X", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "Y", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "Z", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "AA", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "AB", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "AC", last+2), digest.CGFooter[i])
	i++
	f.SetCellValue(sheet, fmt.Sprintf(`%v%d`, "AD", last+2), digest.CGFooter[i])

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	buf, err := f.WriteToBuffer()
	if err != nil {
		report.Error(r, err)
		return nil, err
	}
	return buf.Bytes(), nil

}
