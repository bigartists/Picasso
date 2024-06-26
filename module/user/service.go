package user

import (
	"fmt"
	"picasso/pkg/result"
)

var ServiceGetter Service

func init() {
	ServiceGetter = NewUserGetterImpl()
}

func NewUserGetterImpl() *ServiceGetterImpl {
	return &ServiceGetterImpl{}
}

type ServiceGetterImpl struct {
}

func (this *ServiceGetterImpl) SignIn(username string, password string) (*UserDTO, error) {
	user, err := DaoGetter.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	//if user.Password != password {
	//	err = fmt.Errorf("用户名%s或密码错误", username)
	//	return nil, err
	//}
	// 校验密码
	if !user.CheckPassword(password) {
		err = fmt.Errorf("用户名%s或密码错误", username)
		return nil, err
	}
	userdto := ConvertUserToDTO(user)
	return userdto, nil
}

func (this *ServiceGetterImpl) SignUp(email string, username string, password string) error {
	//return dao.DaoGetter.CreateUser(user)

	if _, err := DaoGetter.FindUserByUsername(username); err != nil {
		return fmt.Errorf("用户名%s已存在", username)
	}
	if _, err := DaoGetter.FindUserByEmail(email); err != nil {
		return fmt.Errorf("邮箱%s已存在", email)
	}

	user := NewModel(WithEmail(email), WithUsername(username), WithPassword(password))
	err := user.GeneratePassword()
	if err != nil {
		return fmt.Errorf("密码加密失败")
	}
	err = DaoGetter.CreateUser(user)

	if err != nil {
		return fmt.Errorf("用户注册失败")
	}

	return nil
}

func (this *ServiceGetterImpl) GetUserList() []*UserDTO {
	users := DaoGetter.FindUserAll()
	userdtos := make([]*UserDTO, len(users))
	for i, user := range users {
		userdtos[i] = ConvertUserToDTO(user)
	}
	return userdtos
}

func (this *ServiceGetterImpl) GetUserDetail(id int64) *result.ErrorResult {
	//TODO implement me
	user := NewModel()
	_, err := DaoGetter.FindUserById(id, user)
	if err != nil {
		return result.Result(nil, err)
	}
	return result.Result(user, nil)
}

func (this *ServiceGetterImpl) CreateUser(user *User) *result.ErrorResult {
	//TODO implement me
	panic("implement me")
}

func (this *ServiceGetterImpl) UpdateUser(id int, user *User) *result.ErrorResult {
	//TODO implement me
	panic("implement me")
}

func (this *ServiceGetterImpl) DeleteUser(id int) *result.ErrorResult {
	//TODO implement me
	panic("implement me")
}

//
//// 创建用户
//func (this *GetterImpl) CreateUser(user *User.UserModelImpl) *result.ErrorResult {
//	db := dbs.Orm.Create(user)
//	if db.Error != nil {
//		return result.Result(nil, db.Error)
//	}
//	return result.Result(user, nil)
//}
//

//
//// 更新用户
//func (this *GetterImpl) UpdateUser(id int, user *User.UserModelImpl) *result.ErrorResult {
//	db := dbs.Orm.Where("id=?", id).Updates(user)
//	if db.Error != nil {
//		return result.Result(nil, db.Error)
//	}
//	return result.Result(user, nil)
//}
//
//// 删除用户
//func (this *GetterImpl) DeleteUser(id int) *result.ErrorResult {
//	user := User.New()
//	db := dbs.Orm.Where("id=?", id).Delete(user)
//	if db.Error != nil {
//		return result.Result(nil, db.Error)
//	}
//	return result.Result(user, nil)
//}
