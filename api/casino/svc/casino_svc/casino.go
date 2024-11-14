package casino_svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type CasinoSvc struct {
	RedisCli   *redis.Redis
	GlobalData *GloablData
}

func NewCasinoSvc(redis *redis.Redis) *CasinoSvc {
	return &CasinoSvc{
		RedisCli:   redis,
		GlobalData: &GloablData{},
	}
}

// CalculateBonus 计算分红
func (s *CasinoSvc) CalculateBonus(ctx context.Context, unBonusRound int64) error {
	// 计算分红
	bonusAmount := float64(0)
	if s.GlobalData.DepositAmount > 0 {
		bonusAmount = float64(s.GlobalData.LastIncome) / float64(s.GlobalData.DepositAmount)
	}

	// 存储分红
	bonus := &Bonus{
		BeginRound:    s.GlobalData.CompletedRound + 1,
		BonusAmount:   bonusAmount,
		ContainRounds: unBonusRound,
	}
	err := s.SaveBonus(ctx, bonus)
	if err != nil {
		logx.Error("save bonus error", logx.Field("error", err))
		return err
	}

	return nil
}

// GetNextRoundUserListRedisKey 获取新增的下一轮用户列表的redis key
func (s *CasinoSvc) GetNextRoundUserListRedisKey() string {
	return "casino:nextRoundUserList"
}

// AddUserToNextRoundUserList 添加用户到新增的下一轮用户列表
func (s *CasinoSvc) AddUserToNextRoundUserList(ctx context.Context, userData *UserData) error {
	userDataBytes, err := json.Marshal(userData)
	if err != nil {
		logx.WithContext(ctx).Errorf("json marshal error. err:%v data:%v", err, userData)
		return err
	}
	err = s.RedisCli.HsetCtx(ctx, s.GetNextRoundUserListRedisKey(), userData.Address, string(userDataBytes))
	if err != nil {
		logx.WithContext(ctx).Errorf("set redis error. err:%v key:%v data:%v", err, userData.Address, string(userDataBytes))
		return err
	}
	return nil
}

// GetUserFromNextRoundUserList 从下一轮用户列表中获取用户
func (s *CasinoSvc) GetUserFromNextRoundUserList(ctx context.Context, address string) (*UserData, error) {
	userDataStr, err := s.RedisCli.HgetCtx(ctx, s.GetNextRoundUserListRedisKey(), address)
	if err != nil && !errors.Is(err, redis.Nil) {
		logx.WithContext(ctx).Errorf("get redis error. err:%v key:%v", err, address)
		return nil, err
	}
	if userDataStr == "" {
		return nil, nil
	}
	userData := &UserData{}
	err = json.Unmarshal([]byte(userDataStr), userData)
	if err != nil {
		logx.WithContext(ctx).Errorf("unmarshal error. err:%v data:%v", err, userDataStr)
		return nil, err
	}
	return userData, nil
}

// GetAllNextRoundUsers 获取下一轮所有的用户
func (s *CasinoSvc) GetAllNextRoundUsers(ctx context.Context) ([]*UserData, error) {
	userListMap, err := s.RedisCli.HgetallCtx(ctx, s.GetNextRoundUserListRedisKey())
	if err != nil {
		logx.WithContext(ctx).Errorf("get users from hashmap error. err:%v", err)
		return nil, err
	}
	userList := make([]*UserData, 0)
	for _, v := range userListMap {
		tmpUserData := &UserData{}
		err = json.Unmarshal([]byte(v), tmpUserData)
		if err != nil {
			logx.WithContext(ctx).Errorf("unmarshal user data error. err:%v data:%v", err, v)
			return nil, err
		}
		userList = append(userList, tmpUserData)
	}
	return userList, nil
}

// ResetNextRoundUserList 重置下一轮用户列表
func (s *CasinoSvc) ResetNextRoundUserList(ctx context.Context) error {
	_, err := s.RedisCli.DelCtx(ctx, s.GetNextRoundUserListRedisKey())
	if err != nil {
		logx.WithContext(ctx).Errorf("del redis key error. err:%v key:%v", err, s.GetNextRoundUserListRedisKey())
		return err
	}
	return nil
}

// GetWithdrawUserListRedisKey 获取提款用户列表的redis key
func (s *CasinoSvc) GetWithdrawUserListRedisKey() string {
	return fmt.Sprintf("casino:withdrawUserList")
}

// GetAllWithdrawUsers 获取所有提款用户
func (s *CasinoSvc) GetAllWithdrawUsers(ctx context.Context) ([]*UserData, error) {
	userListMap, err := s.RedisCli.HgetallCtx(ctx, s.GetWithdrawUserListRedisKey())
	if err != nil {
		logx.WithContext(ctx).Errorf("get users from hashmap error. err:%v", err)
		return nil, err
	}
	userList := make([]*UserData, 0)
	for _, v := range userListMap {
		tmpUserData := &UserData{}
		err = json.Unmarshal([]byte(v), tmpUserData)
		if err != nil {
			logx.WithContext(ctx).Errorf("unmarshal user data error. err:%v data:%v", err, v)
			return nil, err
		}
		userList = append(userList, tmpUserData)
	}
	return userList, nil
}

// HandleNewDeposit 处理新增质押
func (s *CasinoSvc) HandleNewDeposit(ctx context.Context) error {
	// 获取所有新增用户
	userList, err := s.GetAllNextRoundUsers(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("get all next round users error. err:%v", err)
		return err
	}
	// 将新增用户添加到用户列表
	incrAmount := int64(0)
	for _, v := range userList {
		incrAmount += v.DepositAmount
		oldUser, err := s.GetUserData(ctx, v.Address)
		if err != nil {
			logx.WithContext(ctx).Errorf("get user data error. err:%v data:%v", err, v)
			return err
		}
		if oldUser != nil {
			oldUser.DepositAmount += v.DepositAmount
			err = s.SaveUserData(ctx, oldUser)
			if err != nil {
				logx.WithContext(ctx).Errorf("save user data error. err:%v data:%v", err, oldUser)
				return err
			}
		} else {
			err = s.SaveUserData(ctx, v)
			if err != nil {
				logx.WithContext(ctx).Errorf("save user data error. err:%v data:%v", err, v)
				return err
			}
		}
	}
	// 计算新增质押并加入到总质押
	s.GlobalData.DepositAmount += incrAmount
	// 重置新增用户列表
	err = s.ResetNextRoundUserList(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("reset next round user list error. err:%v", err)
		return err
	}
	return nil
}

// ResetWithdrawUserList 重置提款用户列表
func (s *CasinoSvc) ResetWithdrawUserList(ctx context.Context) error {
	_, err := s.RedisCli.DelCtx(ctx, s.GetWithdrawUserListRedisKey())
	if err != nil {
		logx.WithContext(ctx).Errorf("del redis key error. err:%v key:%v", err, s.GetWithdrawUserListRedisKey())
		return err
	}
	return nil
}

// AddUserToWithdrawUserList 添加用户到提现用户列表
func (s *CasinoSvc) AddUserToWithdrawUserList(ctx context.Context, userData *UserData) error {
	if userData == nil {
		logx.WithContext(ctx).Errorf("empty user data")
		return fmt.Errorf("empty user data")
	}
	userDataBytes, err := json.Marshal(userData)
	if err != nil {
		logx.WithContext(ctx).Errorf("json marshal error. err:%v data:%v", err, userData)
		return err
	}
	err = s.RedisCli.HsetCtx(ctx, s.GetWithdrawUserListRedisKey(), userData.Address, string(userDataBytes))
	if err != nil {
		logx.WithContext(ctx).Errorf("set redis error. err:%v key:%v data:%v", err, userData.Address, string(userDataBytes))
		return err
	}
	return nil
}

// GetUserFromWithdrawUserList 从提现用户列表中获取用户
func (s *CasinoSvc) GetUserFromWithdrawUserList(ctx context.Context, address string) (*UserData, error) {
	userDataStr, err := s.RedisCli.HgetCtx(ctx, s.GetWithdrawUserListRedisKey(), address)
	if err != nil && !errors.Is(err, redis.Nil) {
		logx.WithContext(ctx).Errorf("get redis error. err:%v key:%v", err, address)
		return nil, err
	}
	if userDataStr == "" {
		return nil, nil
	}
	userData := &UserData{}
	err = json.Unmarshal([]byte(userDataStr), userData)
	if err != nil {
		logx.WithContext(ctx).Errorf("unmarshal user data error. err:%v data:%v", err, userDataStr)
		return nil, err
	}
	return userData, nil
}

// DeleteUserData 删除用户数据
func (s *CasinoSvc) DeleteUserData(ctx context.Context, address string) error {
	_, err := s.RedisCli.HdelCtx(ctx, s.GetUserListRedisKey(), address)
	if err != nil {
		logx.WithContext(ctx).Errorf("del redis key error. err:%v key:%v", err, address)
		return err
	}
	return nil
}

// HandleWithdraw 处理提款
func (s *CasinoSvc) HandleWithdraw(ctx context.Context) error {
	// 获取所有提款用户
	userList, err := s.GetAllWithdrawUsers(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("get all withdraw users error. err:%v", err)
		return err
	}

	decrAmount := int64(0)
	// 将提款用户从用户列表中删除并计算提款金额
	for _, v := range userList {
		tmpUserData, err := s.GetUserData(ctx, v.Address)
		if err != nil {
			logx.WithContext(ctx).Errorf("get user data error. err:%v data:%v", err, v)
			return err
		}
		decrAmount += tmpUserData.DepositAmount

		err = s.DeleteUserData(ctx, v.Address)
		if err != nil {
			logx.WithContext(ctx).Errorf("delete user data error. err:%v data:%v", err, v)
			return err
		}
	}
	// 重置提款用户列表
	err = s.ResetWithdrawUserList(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("reset withdraw user list error. err:%v", err)
		return err
	}
	return nil
}

// NextRound 开启新的一轮
func (s *CasinoSvc) NextRound(ctx context.Context, unBonusRound int64) error {
	// 处理分红
	err := s.CalculateBonus(ctx, unBonusRound)
	if err != nil {
		logx.WithContext(ctx).Errorf("calculate bonus error: %v", err)
		return err
	}

	// 处理新增质押
	err = s.HandleNewDeposit(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("handle new deposit error: %v", err)
		return err
	}

	// 处理取款
	err = s.HandleWithdraw(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("handle withdraw error: %v", err)
		return err
	}

	// 更新全局数据
	s.GlobalData.CompletedRound += unBonusRound
	err = s.SaveGlobalData(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("save global data error: %v", err)
		return err
	}
	return nil
}

// LoadGlobalData 加载全局数据
func (s *CasinoSvc) LoadGlobalData(ctx context.Context) error {
	globalDataStr, err := s.RedisCli.GetCtx(ctx, s.GetGlobalDataRedisKey())
	if err != nil {
		logx.WithContext(ctx).Errorf("get global data error: %v", err)
		return err
	}
	if globalDataStr == "" {
		s.GlobalData = &GloablData{}
		return nil
	}
	globalData := &GloablData{}
	err = json.Unmarshal([]byte(globalDataStr), globalData)
	if err != nil {
		logx.WithContext(ctx).Errorf("unmarshal global data error: %v", err)
		return err
	}
	s.GlobalData = globalData
	return nil
}

// SaveGlobalData 保存全局数据
func (s *CasinoSvc) SaveGlobalData(ctx context.Context) error {
	globalDataBytes, err := json.Marshal(s.GlobalData)
	if err != nil {
		logx.WithContext(ctx).Errorf("json marshal error. err:%v data:%v", err, s.GlobalData)
		return err
	}
	err = s.RedisCli.SetCtx(ctx, s.GetGlobalDataRedisKey(), string(globalDataBytes))
	if err != nil {
		logx.WithContext(ctx).Errorf("set redis error. err:%v key:%v data:%v", err, s.GetGlobalDataRedisKey(), string(globalDataBytes))
		return err
	}
	return nil
}

// GetGlobalDataRedisKey 获取全局数据redis key
func (s *CasinoSvc) GetGlobalDataRedisKey() string {
	return "casino:globalData"
}

// GetUserListRedisKey 获取用户列表redis key
func (s *CasinoSvc) GetUserListRedisKey() string {
	return fmt.Sprintf("casino:userList")
}

// GetUserData 获取用户数据
func (s *CasinoSvc) GetUserData(ctx context.Context, address string) (userData *UserData, err error) {
	userDataStr, err := s.RedisCli.HgetCtx(ctx, s.GetUserListRedisKey(), address)
	if err != nil && !errors.Is(err, redis.Nil) {
		logx.WithContext(ctx).Errorf("get global data error: %v", err)
		return nil, err
	}
	if userDataStr == "" {
		return nil, nil
	}
	userData = &UserData{}
	err = json.Unmarshal([]byte(userDataStr), userData)
	if err != nil {
		logx.WithContext(ctx).Errorf("unmarshal global data error: %v", err)
		return nil, err
	}
	return userData, nil
}

// SaveUserData 保存用户数据
func (s *CasinoSvc) SaveUserData(ctx context.Context, userData *UserData) error {
	if userData == nil {
		logx.WithContext(ctx).Errorf("empty user data")
		return fmt.Errorf("empty user data")
	}
	userDataBytes, err := json.Marshal(userData)
	if err != nil {
		logx.WithContext(ctx).Errorf("json marshal error. err:%v data:%v", err, userData)
		return err
	}
	err = s.RedisCli.HsetCtx(ctx, s.GetUserListRedisKey(), userData.Address, string(userDataBytes))
	if err != nil {
		logx.WithContext(ctx).Errorf("set redis error. err:%v key:%v data:%v", err, userData.Address, string(userDataBytes))
		return err
	}
	return nil
}

// GetBonusListRedisKey 获取分红列表的redis key
func (s *CasinoSvc) GetBonusListRedisKey() string {
	return fmt.Sprintf("casino:bonusList")
}

// GetBonus 获取分红
func (s *CasinoSvc) GetBonus(ctx context.Context, round int64) (*Bonus, error) {
	bonusStr, err := s.RedisCli.HgetCtx(ctx, s.GetBonusListRedisKey(), fmt.Sprintf("%v", round))
	if err != nil {
		logx.WithContext(ctx).Error("get bonus error. err:%v round:%v", err, round)
		return nil, err
	}
	bonus := &Bonus{}
	err = json.Unmarshal([]byte(bonusStr), bonus)
	if err != nil {
		logx.WithContext(ctx).Error("unmarshal bonus error. err:%v round:%v", err, round)
		return nil, err
	}
	return bonus, nil
}

// SaveBonus 保存分红
func (s *CasinoSvc) SaveBonus(ctx context.Context, bonus *Bonus) error {
	if bonus == nil {
		logx.WithContext(ctx).Errorf("empty bonus")
		return fmt.Errorf("empty bonus")
	}
	bonusBytes, err := json.Marshal(bonus)
	if err != nil {
		logx.WithContext(ctx).Error("json marshal error. err:%v data:%v", err, bonus)
		return err
	}
	for i := int64(0); i < bonus.ContainRounds; i++ {
		err = s.RedisCli.HsetCtx(ctx, s.GetBonusListRedisKey(), fmt.Sprintf("%v", bonus.BeginRound+i), string(bonusBytes))
		if err != nil {
			logx.WithContext(ctx).Error("set redis error. err:%v key:%v data:%v", err, fmt.Sprintf("%v", bonus.BeginRound), string(bonusBytes))
			return err
		}
	}
	return nil
}

// GetCurrentRoundAndUnBonusRound 获取当前轮次和未分红轮次以及是否可以开启下一轮
func (s *CasinoSvc) GetCurrentRoundAndUnBonusRound(ctx context.Context, blockSeq int64) (curRound int64, unBonusRound int64, canNextRound bool, err error) {
	// 当前轮次
	remainder := blockSeq % OneDayBlocks
	curRound = int64(0)
	if remainder == 0 {
		curRound = blockSeq / OneDayBlocks
	} else {
		curRound = blockSeq/OneDayBlocks + 1
	}

	// 未分红轮次
	canNextRound = false
	unBonusRound = curRound - s.GlobalData.CompletedRound - 1
	if unBonusRound >= 1 {
		canNextRound = true
	}
	return curRound, unBonusRound, canNextRound, nil
}

// StartNextRound 开启下一轮
func (s *CasinoSvc) StartNextRound(ctx context.Context, unBonusRound int64) error {
	// 计算分红
	err := s.CalculateBonus(ctx, unBonusRound)
	if err != nil {
		logx.WithContext(ctx).Errorf("calculate bonus error: %v", err)
		return err
	}

	// 处理取款
	err = s.HandleWithdraw(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("handle withdraw error: %v", err)
		return err
	}

	// 处理新增质押
	err = s.HandleNewDeposit(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("handle new deposit error: %v", err)
		return err
	}

	s.GlobalData.LastIncome = s.GlobalData.PresentIncome
	s.GlobalData.PresentIncome = 0
	s.GlobalData.CompletedRound += unBonusRound

	err = s.SaveGlobalData(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("save global data error: %v", err)
		return err
	}
	return nil
}
