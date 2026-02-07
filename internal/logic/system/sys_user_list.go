package system

import (
	"context"
	"sagooiot/internal/dao"
	"sagooiot/internal/model"
	"sagooiot/internal/model/entity"
	"sagooiot/internal/service"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// UserList 用户列表
func (s *sSysUser) UserList(ctx context.Context, input *model.UserListDoInput) (total int, out []*model.UserListOut, err error) {
	if input == nil {
		input = &model.UserListDoInput{}
	}
	input.PaginationInput = model.EnsurePaginationInput(input.PaginationInput)

	m := dao.SysUser.Ctx(ctx)
	if input.KeyWords != "" {
		keyWords := "%" + input.KeyWords + "%"
		m = m.Where("user_name like ? or  user_nickname like ?", keyWords, keyWords)
	}
	if input.DeptId != 0 {
		//m = m.Where(dao.SysUser.Columns().DeptId, req.DeptId)
		deptIds, _ := s.getSearchDeptIds(ctx, gconv.Int64(input.DeptId))
		m = m.WhereIn(dao.SysUser.Columns().DeptId, deptIds)
	}
	if input.UserName != "" {
		m = m.Where(dao.SysUser.Columns().UserName, input.UserName)
	}
	if input.Status != -1 {
		m = m.Where(dao.SysUser.Columns().Status, input.Status)
	}
	if input.Mobile != "" {
		m = m.WhereLike(dao.SysUser.Columns().Mobile, "%"+input.Mobile+"%")
	}
	if len(input.DateRange) > 0 {
		m = m.Where("created_at >=? AND created_at <?", input.DateRange[0], gtime.NewFromStrFormat(input.DateRange[1], "Y-m-d").AddDate(0, 0, 1))
	}
	m = m.Where(dao.SysUser.Columns().IsDeleted, 0)

	//获取总数
	total, err = m.Count()
	if err != nil {
		err = gerror.New("获取数据失败")
		return
	}
	//获取用户信息
	err = m.Page(input.PageNum, input.PageSize).OrderDesc(dao.SysUser.Columns().CreatedAt).Scan(&out)
	if err != nil {
		err = gerror.New("获取用户列表失败")
	}

	deptIds := g.Slice{}
	for _, v := range out {
		deptIds = append(deptIds, v.DeptId)
	}
	deptData, _ := GetDeptNameDict(ctx, deptIds)
	for _, v := range out {
		err = gconv.Scan(deptData[v.DeptId], &v.Dept)
		v.RolesNames = getUserRoleName(ctx, gconv.Int(v.Id))
	}
	return
}

func getUserRoleName(ctx context.Context, userId int) (rolesNames string) {

	rolesData, _ := gcache.Get(ctx, "RoleListAtName"+gconv.String(userId))
	if rolesData != nil {
		rolesNames = rolesData.String()
		return
	}

	var roleIds []int
	//查询用户角色
	userRoleInfo, _ := service.SysUserRole().GetInfoByUserId(ctx, userId)
	if userRoleInfo != nil {
		for _, userRole := range userRoleInfo {
			roleIds = append(roleIds, userRole.RoleId)
		}
	}
	if roleIds != nil {
		//获取所有的角色信息
		roleInfo, _ := service.SysRole().GetInfoByIds(ctx, roleIds)
		var roleNames []string
		if roleInfo != nil {
			for _, role := range roleInfo {
				roleNames = append(roleNames, role.Name)
			}
			rolesNames = strings.Join(roleNames, ",")

			//放入缓存
			effectiveTime := time.Minute * 5
			_, _ = gcache.SetIfNotExist(ctx, "RoleListAtName"+gconv.String(userId), rolesNames, effectiveTime)

			return

		}
	}

	return ""
}

func GetDeptNameDict(ctx context.Context, ids g.Slice) (dict map[int64]model.DetailDeptRes, err error) {
	var depts []model.DetailDeptRes
	dict = make(map[int64]model.DetailDeptRes)
	err = dao.SysDept.Ctx(ctx).Where(dao.SysDept.Columns().DeptId, ids).Scan(&depts)
	for _, d := range depts {
		dict[d.DeptId] = d
	}
	return
}

// GetUserByIds 根据ID数据获取用户信息
func (s *sSysUser) GetUserByIds(ctx context.Context, id []int) (data []*entity.SysUser, err error) {
	m := dao.SysUser.Ctx(ctx)

	err = m.Where(g.Map{
		dao.SysUser.Columns().IsDeleted: 0,
	}).WhereIn(dao.SysUser.Columns().Id, id).Scan(&data)
	return
}

// GetAll 获取所有用户信息
func (s *sSysUser) GetAll(ctx context.Context) (data []*entity.SysUser, err error) {
	m := dao.SysUser.Ctx(ctx)

	err = m.Where(g.Map{
		dao.SysUser.Columns().IsDeleted: 0,
	}).Scan(&data)
	return
}

// GetUsersByCodes 根据用户编码数组批量获取用户信息
func (s *sSysUser) GetUsersByCodes(ctx context.Context, codes []string) (data []*entity.SysUser, err error) {
	err = dao.SysUser.Ctx(ctx).
		WhereIn(dao.SysUser.Columns().Code, codes).
		Where(dao.SysUser.Columns().IsDeleted, 0).
		Scan(&data)
	return
}
