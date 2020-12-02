# 个人理解
sql.ErrNoRows这个错误应该在dao层中处理掉，不应该在service层去做判断条件，这样减少功能耦合。
我的处理方式是在dao层若发现sql.ErrNoRows这个错误便返回空数据和nil error， 这样的处理方式是不是错的？

# 代码实现

## model层
`
type Persion struct {
	Name string
	Age  string
}

func (p *Persion) TableName() string {
	return "persion"
}
`

## dao层
`
type DBService struct {
	db *Db
}

func (this *DBService) FindPersion() ([]model.Persion, error) {
	data, err := this.db.query("select * from persion")
	if err != nil {
		if errors.Is(sql.ErrNoRows) {
			// 当sql.ErrNoRows时，这里把error吞掉可不可以，直接返回一个nil
			return nil, nil
		}
		return nil, errors.Wrap(err, "query persion error: ")
	}
	return data, err
}
`
## service层
`
type Service struct {
	db *dao.DBService
}

func NewService() *Service {
	return &Service{db: dao.NewDBService()}
}

func (this *Service) FindPersion() ([]model.Persion, error) {
	return this.db.FindPersion()
}
`

## main
`
func main() {
	service := service.NewService()
	data, err := service.FindPersion()
	if err != nil {
		fmt.Println("find persion error: %+v", err)
	}
	fmt.Println("persion info: %+v", data)
}
`
