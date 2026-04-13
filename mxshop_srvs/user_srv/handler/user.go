package handler

import (
	"context"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

// Paginate 分页函数
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// ModelToResponse 将 User 模型转换为 UserInfo 响应
func ModelToResponse(user *model.User) *proto.UserInfo {
	rsp := &proto.UserInfo{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		rsp.Birthday = uint64(user.Birthday.Unix())
	}
	return rsp
}

// 加密密码（生成哈希）
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10，可根据需要调整（4-31）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// 验证密码（比对明文和哈希，不需要解密）
func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

type UserServer struct{
	proto.UnimplementedUserServer
}

// GetUserList 获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(&user)
		rsp.Data = append(rsp.Data, userInfoRsp)
	}

	return rsp, nil
}

// GetUserByMobile 根据手机号获取用户信息
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfo, error) {
	var user model.User
	result := global.DB.Where("mobile = ?", req.Mobile).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(&user), nil
}

// GetUserById 根据用户ID获取用户信息
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfo, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(&user), nil
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserInfo, error) {
	// 查询用户是否存在
	var user model.User
	result := global.DB.Where("mobile = ?", req.Mobile).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	// 密码加密
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "系统内部错误")
	}
	user.Password = hashedPassword

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return ModelToResponse(&user), nil
}

// UpdateUser 更新用户信息
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthday := time.Unix(int64(req.Birthday), 0)
	user.NickName = req.NickName
	user.Birthday = &birthday
	user.Gender = req.Gender

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &emptypb.Empty{}, nil
}

// CheckPassword 验证密码
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordRequest) (*proto.CheckPasswordResponse, error) {
	CheckPassword(req.EncryptedPassword, req.Password)
	if err := CheckPassword(req.EncryptedPassword, req.Password); err != nil {
		return &proto.CheckPasswordResponse{Success: false}, nil
	}

	return &proto.CheckPasswordResponse{Success: true}, nil
}


