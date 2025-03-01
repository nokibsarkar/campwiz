package types

type IDType string

func (i *IDType) String() string {
	return string(*i)
}
func (i *IDType) GormDataType() string {
	return "varchar(255)"
}
