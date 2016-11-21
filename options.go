package goSmartSheet

import "github.com/pkg/errors"

//RowPostOptions is used in conjuction with Adding a Row to a sheet to specific location
type RowPostOptions int16

const (
	//ToTop will Add or move the Row to the top of the Sheet.
	ToTop RowPostOptions = iota
	//ToBottom will Add or move the Row to the bottom of the Sheet.
	ToBottom
	//Above will Add or move the Row directly above the specified sibling Row (at the same hierarchical level).
	//Sibling Row must be populated for this option to work
	Above
)

//PostOptions is used during a post to control the level of validation / adjustments performed by the client
type PostOptions int16

const (
	//NormalValidation performs the default logic for that specific operation
	NormalValidation PostOptions = 1 << iota
	//IgnoreColumnLengthValidation will disable all column validation and assueme the calling application build the columns / rows correctly
	IgnoreColumnLengthValidation
	//IgnoreRightMostColumns will fix / adjust the leading columns and then ignore the rest of the columns provided
	IgnoreRightMostColumns
)

//ValidateCellsInRow will validate that the cells match the columns within the sheet based on the specified PostOptions
func ValidateCellsInRow(cells []Cell, sheetCols []Column, opt PostOptions) error {
	switch opt {
	case NormalValidation:
		if len(sheetCols) != len(cells) {
			return errors.New("Cells within a row  must match columns in sheet")
		}
	case IgnoreRightMostColumns:
		//only validate that it does not have more columns
		if len(sheetCols) < len(cells) {
			return errors.New("Cells within a row cannot be greater than the columns within the sheet")
		}
	}

	return nil
}
