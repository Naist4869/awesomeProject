package mongo

import (
	"context"
	"reflect"

	"github.com/Naist4869/awesomeProject/model/usermodel"

	"github.com/Naist4869/awesomeProject/dataservice"

	"go.uber.org/zap"

	"github.com/Naist4869/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DBUserVersion = 1
	DBUserKey     = "users"
)

var (
	defaultSort = []string{"_id"}
)

type UserClient struct {
	Client      *mongo.Client
	DB          string                       // 数据库名
	collections map[string]*mongo.Collection // 数据表map
	Logger      *log.Logger
}

func NewUserClient(client *mongo.Client, DB string, logger *log.Logger) *UserClient {
	return &UserClient{Client: client, DB: DB, collections: make(map[string]*mongo.Collection, 3), Logger: logger}
}

func (uc *UserClient) Keys() map[string]*Spec {
	specs := make(map[string]*Spec, 1)
	userSpec, err := NewSpec(DBUserVersion, func() interface{} {
		return &usermodel.User{}
	}, func(data interface{}) error {
		return nil
	})
	if err != nil {
		uc.Logger.Fatal("构建updater失败:", zap.Error(err))
	}
	specs[DBUserKey] = userSpec
	return specs
}
func (uc *UserClient) Init() error {
	keys := uc.Keys()
	for key, spec := range keys {
		collection := uc.Client.Database(uc.DB).Collection(key)
		uc.collections[key] = collection
		if spec != nil {
			spec.SetCollection(collection)
		}
	}
	updater, err := NewUpdater(keys, uc.Logger)
	if err != nil {
		return errors.Wrap(err, "构建collection升级器失败")
	}
	if err := updater.Update(); err != nil {
		return errors.Wrap(err, "检查并升级collection数据失败")
	}
	return nil
}

func (uc *UserClient) Remove(userID int64) (err error) {
	panic("implement me")
}

func (uc *UserClient) Insert(u *usermodel.User) (err error) {
	if _, err = uc.collections[DBUserKey].InsertOne(context.Background(), *u); err != nil {
		return
	}
	return
}

func (uc *UserClient) FindByID(ctx context.Context, userID int64) (user *usermodel.User, err error) {
	query := bson.M{
		usermodel.IDField:     userID,
		usermodel.StatusField: usermodel.Normal,
	}
	users, _, err := uc.queryUser(ctx, query, nil, 0, 1, nil, nil)
	if err != nil {
		return nil, errors.New("根据id获取用户错误")
	}
	if len(users) == 0 {
		return nil, dataservice.NewErrIDNotFound(userID)
	}
	return users[0], nil
}

func (uc *UserClient) FindByIDs(userIDs []int64) (users []*usermodel.User, err error) {
	panic("implement me")
}

func (uc *UserClient) FindByPID(userID int64) (user *usermodel.User, count int64, err error) {
	panic("implement me")
}

func (uc *UserClient) Update(user *usermodel.User) (err error) {
	panic("implement me")
}
func (uc *UserClient) FindByPhone(phone string) (*usermodel.User, error) {
	query := bson.M{
		usermodel.PhoneField:  phone,
		usermodel.StatusField: usermodel.Normal,
	}
	users, _, err := uc.queryUser(context.Background(), query, nil, 0, 1, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "根据手机号获取用户错误")
	}
	if len(users) == 0 {
		return nil, dataservice.NewErrPhoneNotFound(phone)
	}
	return users[0], nil

}
func (uc *UserClient) queryUser(ctx context.Context, query bson.M, sort []string, start, limit int64, include, exclude []string, collations ...*options.Collation) (records []*usermodel.User, count int64, err error) {
	var (
		data    = &usermodel.User{}
		results []interface{}
	)

	if len(sort) == 0 {
		sort = defaultSort
	}
	if len(query) == 0 {
		query = make(bson.M, 1)
	}
	query[usermodel.DeletedField] = false // 查询未删除的数据
	results, count, err = uc.baseQuery(uc.collections[DBUserKey], ctx, query, sort, start, limit, include, exclude, data, collations...)
	if err != nil {
		return
	}
	records = make([]*usermodel.User, 0, len(results))
	for _, data := range results {
		records = append(records, data.(*usermodel.User))
	}
	return
}

func (uc *UserClient) baseQuery(collection *mongo.Collection, ctx context.Context, query bson.M, sort []string, start, limit int64, include []string, exclude []string, data interface{}, collations ...*options.Collation) (result []interface{}, count int64, err error) {
	var selection bson.M
	if selection, err = MakeSelect(include, exclude); len(selection) == 0 {
		selection = nil
	}
	option := &options.FindOptions{
		Sort:       convertSort(sort),
		Limit:      &limit,
		Skip:       &start,
		Projection: selection,
	}
	switch len(collations) {
	case 0:
	case 1:
		option.Collation = collations[0]
	default:
		err = errors.New("collations参数最多只有1个")
		return
	}
	var cursor *mongo.Cursor
	if cursor, err = collection.Find(ctx, query, option); err != nil {
		err = errors.Wrap(err, "find驱动")
		return
	}
	t := reflect.TypeOf(data)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		element := reflect.New(t)
		if err = cursor.Decode(element.Interface()); err != nil {
			err = errors.Wrap(err, "bson解码")
			return
		}
		result = append(result, element.Elem().Interface())
		count++
	}
	return
}
func MakeSelect(include, exclude []string) (selection bson.M, err error) {
	if err = validateSelect(include, exclude); err != nil {
		return
	}
	selection = makeSelect(include, exclude)
	return
}

func validateSelect(include, exclude []string) (err error) {
	if len(include) != 0 && len(exclude) != 0 {
		err = errors.New("两个参数必须至少有一个为空")
		return
	}
	return
}

func makeSelect(include, exclude []string) (fields bson.M) {
	fields = make(bson.M, 2)
	for _, field := range include {
		fields[field] = 1
	}
	for _, field := range exclude {
		fields[field] = 0
	}
	return
}
func convertSort(sorts []string) (result bson.D) {
	result = make([]bson.E, 0, len(sorts))
	value := 1
	for _, sort := range sorts {
		if sort[:1] == "-" {
			sort = sort[1:]
			value = -1
		}
		result = append(result, bson.E{
			Key:   sort,
			Value: value,
		})
	}
	return
}
