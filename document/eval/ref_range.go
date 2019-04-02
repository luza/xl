package eval

type rangeRef struct {
	Value

	SheetIdx int
	CellFrom Axis
	CellTo   Axis
}

//func (r *cellRef) BoolValue(ec *eval.Context) (bool, error) {
//	//if r.CellTo != nil {
//	//	return false, NewError(ErrorKindCasting, "unable to cast range to bool")
//	//}
//	if ec.Visited(r) {
//		return false, eval.NewError(eval.ErrorKindRef, "circular reference")
//	}
//	ec.AddVisited(r)
//	return ec.LinkRegistry.BoolValue(ec, r.SheetIdx, r.Cell.X, r.Cell.Y)
//}
//
//func (r *cellRef) DecimalValue(ec *eval.Context) (decimal.Decimal, error) {
//	//if r.CellTo != nil {
//	//	return decimal.Zero, NewError(ErrorKindCasting, "unable to cast range to decimal")
//	//}
//	if ec.Visited(r) {
//		return decimal.Zero, eval.NewError(eval.ErrorKindRef, "circular reference")
//	}
//	ec.AddVisited(r)
//	return ec.LinkRegistry.DecimalValue(ec, r.SheetIdx, r.Cell.X, r.Cell.Y)
//}
//
//func (r *cellRef) StringValue(ec *eval.Context) (string, error) {
//	//if r.CellTo != nil {
//	//	return "", NewError(ErrorKindCasting, "unable to cast range to string")
//	//}
//	if ec.Visited(r) {
//		return "", eval.NewError(eval.ErrorKindRef, "circular reference")
//	}
//	ec.AddVisited(r)
//	return ec.LinkRegistry.StringValue(ec, r.SheetIdx, r.Cell.X, r.Cell.Y)
//}
