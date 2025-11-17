package base

import (
	"context"

	pb "github.com/example/aichat/backend/api/base"
	"github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/models/generator/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SystemService struct {
	pb.UnimplementedSystemServer
	menu *base.SysMenuUseCase
}

func NewSystemService(menu *base.SysMenuUseCase) *SystemService {
	return &SystemService{
		menu: menu,
	}
}

func (s *SystemService) Menu(ctx context.Context, req *emptypb.Empty) (*pb.MenuReply, error) {
	menus, err := s.menu.GetRouter(ctx)
	if err != nil {
		return nil, err
	}
	var pbMenus = make([]*pb.Menu, 0, len(menus))
	for _, m := range menus {
		pbMenus = append(pbMenus, &pb.Menu{
			Id:        m.ID,
			ParentID:  m.ParentID,
			Name:      m.Name,
			Path:      m.Path,
			Component: m.Component,
			// Sort:      m.Sort,
			Hidden:     m.Hidden,
			AlwaysShow: m.AlwaysShow,
			Redirect:   m.Redirect,
			Meta: &pb.MenuMeta{
				Title:   m.Meta.Title,
				Icon:    m.Meta.Icon,
				NoCache: m.Meta.NoCache,
			},
		})
	}
	return &pb.MenuReply{
		Menu: pbMenusToTree(pbMenus, 0),
	}, nil
}

func (s *SystemService) AllMenu(ctx context.Context, req *emptypb.Empty) (*pb.AllMenuReply, error) {
	menus, err := s.menu.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var pbMenus = make([]*pb.MenuItem, 0, len(menus))
	for _, m := range menus {
		// 		createBy?: any;
		//   createTime: string;
		//   updateBy?: any;
		//   updateTime?: any;
		//   remark?: any;
		//   menuId: number;
		//   menuName: string;
		//   parentName?: string;
		//   parentId: number;
		//   orderNum: number;
		//   path: string;
		//   component?: string;
		//   query: string;
		//   isFrame: string;
		//   isCache: string;
		//   menuType: string;
		//   visible: string;
		//   status: string;
		//   perms: string;
		//   icon: string;
		pbMenus = append(pbMenus, &pb.MenuItem{
			CreateBy:   "",
			CreateTime: m.CreatedAt.String(),
			UpdateBy:   "",
			UpdateTime: m.UpdatedAt.String(),
			Remark:     "",
			MenuId:     m.ID,
			MenuName:   m.Meta.Title,
			ParentId:   m.ParentID,
			ParentName: "",
			OrderNum:   int32(m.Sort),
			Path:       m.Path,
			Component:  m.Component,
			Query:      "",
			IsFrame:    "1",
			IsCache:    "1",
			MenuType:   m.Type.String(),
			Visible:    "1",
			Status:     "1",
			Perms:      "",
			Icon:       m.Meta.Icon,
		})
	}
	return &pb.AllMenuReply{
		Menus: pbMenus,
	}, nil
}

// pbMenus to tree
func pbMenusToTree(menus []*pb.Menu, parentID int64) []*pb.Menu {
	var tree []*pb.Menu
	for _, menu := range menus {
		if menu.ParentID == parentID {
			// 递归查找子菜单
			menu.Children = pbMenusToTree(menus, menu.Id)
			tree = append(tree, menu)
		}
	}
	return tree
}

// add
func (s *SystemService) AddSysMenu(ctx context.Context, req *pb.AddSysMenuRequest) (*emptypb.Empty, error) {
	t := &model.SysMenu{
		Name:      req.Component,
		Component: req.Component,
		Path:      req.Path,
		// Query:     req.Query,
		// Redirect:  req.Redirect,
		ParentID: req.ParentId,
		Sort:     int(req.OrderNum),
		Type:     model.ToMenuType(req.MenuType),
		Hidden:   req.Visible == "2",
		// AlwaysShow: req.AlwaysShow,
		Meta: &model.Meta{
			Title:   req.MenuName,
			Icon:    req.Icon,
			NoCache: req.IsCache == "2",
		},
		PermsCode: req.Perms,
		Remark:    req.Remark,
	}
	t.New()
	err := s.menu.Add(ctx, t)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateSysMenu
func (s *SystemService) UpdateSysMenu(ctx context.Context, req *pb.AddSysMenuRequest) (*emptypb.Empty, error) {
	t := &model.SysMenu{
		Name:      req.Component,
		Component: req.Component,
		Path:      req.Path,
		// Query:     req.Query,
		// Redirect:  req.Redirect,
		ParentID: req.ParentId,
		Sort:     int(req.OrderNum),
		Type:     model.ToMenuType(req.MenuType),
		Hidden:   req.Visible == "2",
		// AlwaysShow: req.AlwaysShow,
		Meta: &model.Meta{
			Title:   req.MenuName,
			Icon:    req.Icon,
			NoCache: req.IsCache == "2",
		},
		PermsCode: req.Perms,
		Remark:    req.Remark,
	}
	t.ID = req.MenuId
	err := s.menu.Update(ctx, t)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteSysMenu
func (s *SystemService) DeleteSysMenu(ctx context.Context, req *pb.DeleteSysMenuRequest) (*emptypb.Empty, error) {
	err := s.menu.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// GetSysMenu
func (s *SystemService) GetSysMenu(ctx context.Context, req *pb.GetSysMenuRequest) (*pb.MenuItem, error) {
	m, err := s.menu.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.MenuItem{
			CreateBy:   "",
			CreateTime: m.CreatedAt.String(),
			UpdateBy:   "",
			UpdateTime: m.UpdatedAt.String(),
			Remark:     "",
			MenuId:     m.ID,
			MenuName:   m.Meta.Title,
			ParentId:   m.ParentID,
			ParentName: "",
			OrderNum:   int32(m.Sort),
			Path:       m.Path,
			Component:  m.Component,
			Query:      "",
			IsFrame:    "1",
			IsCache:    "1",
			MenuType:   m.Type.String(),
			Visible:    "1",
			Status:     "1",
			Perms:      "",
			Icon:       m.Meta.Icon,
		}, nil
}