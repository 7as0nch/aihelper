package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// WorkflowStep 定义工作流步骤接口
type WorkflowStep interface {
	Execute(ctx context.Context, data interface{}) (interface{}, error)
	GetName() string
}

// DatabaseQueryStep 数据库查询步骤
type DatabaseQueryStep struct {
	name string
	db   *gorm.DB
	log  *log.Helper
}

func NewDatabaseQueryStep(db *gorm.DB, logger log.Logger) *DatabaseQueryStep {
	return &DatabaseQueryStep{
		name: "database_query",
		db:   db,
		log:  log.NewHelper(logger),
	}
}

func (s *DatabaseQueryStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	s.log.Infof("执行数据库查询步骤")
	
	// 这里可以根据传入的数据执行不同的查询
	// 例如：查询用户反馈数据
	var results []map[string]interface{}
	if err := s.db.WithContext(ctx).Table("user_feedbacks").Find(&results).Error; err != nil {
		s.log.Errorf("数据库查询失败: %v", err)
		return nil, err
	}
	
	s.log.Infof("数据库查询完成，获取到 %d 条记录", len(results))
	return results, nil
}

func (s *DatabaseQueryStep) GetName() string {
	return s.name
}

// ExcelExportStep Excel导出步骤
type ExcelExportStep struct {
	name string
	log  *log.Helper
}

func NewExcelExportStep(logger log.Logger) *ExcelExportStep {
	return &ExcelExportStep{
		name: "excel_export",
		log:  log.NewHelper(logger),
	}
}

func (s *ExcelExportStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	s.log.Infof("执行Excel导出步骤")
	
	// 创建新的Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			s.log.Errorf("关闭Excel文件失败: %v", err)
		}
	}()
	
	// 设置工作表名称
	index, err := f.NewSheet("数据导出")
	if err != nil {
		s.log.Errorf("创建工作表失败: %v", err)
		return nil, err
	}
	
	// 设置表头
	headers := []string{"ID", "用户ID", "反馈内容", "评分", "创建时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("数据导出", cell, header)
	}
	
	// 填充数据
	if data != nil {
		if records, ok := data.([]map[string]interface{}); ok {
			for row, record := range records {
				col := 1
				cell, _ := excelize.CoordinatesToCellName(col, row+2)
				f.SetCellValue("数据导出", cell, record["id"])
				
				col++
				cell, _ = excelize.CoordinatesToCellName(col, row+2)
				f.SetCellValue("数据导出", cell, record["user_id"])
				
				col++
				cell, _ = excelize.CoordinatesToCellName(col, row+2)
				f.SetCellValue("数据导出", cell, record["content"])
				
				col++
				cell, _ = excelize.CoordinatesToCellName(col, row+2)
				f.SetCellValue("数据导出", cell, record["rating"])
				
				col++
				cell, _ = excelize.CoordinatesToCellName(col, row+2)
				f.SetCellValue("数据导出", cell, record["created_at"])
			}
		}
	}
	
	// 设置活动工作表
	f.SetActiveSheet(index)
	
	// 生成文件名
	filename := fmt.Sprintf("export_%s.xlsx", time.Now().Format("20060102_150405"))
	
	// 保存文件
	if err := f.SaveAs(filename); err != nil {
		s.log.Errorf("保存Excel文件失败: %v", err)
		return nil, err
	}
	
	s.log.Infof("Excel导出完成，文件名: %s", filename)
	return filename, nil
}

func (s *ExcelExportStep) GetName() string {
	return s.name
}

// WorkflowStatus 工作流状态
type WorkflowStatus string

const (
	WorkflowPending   WorkflowStatus = "pending"
	WorkflowRunning   WorkflowStatus = "running"
	WorkflowCompleted WorkflowStatus = "completed"
	WorkflowFailed    WorkflowStatus = "failed"
)

// WorkflowExecution 工作流执行结果
type WorkflowExecution struct {
	ID         string         `json:"id"`
	Status     WorkflowStatus `json:"status"`
	Steps      []string       `json:"steps"`
	Results    []interface{}  `json:"results"`
	Error      string         `json:"error,omitempty"`
	StartedAt  time.Time      `json:"started_at"`
	FinishedAt *time.Time     `json:"finished_at,omitempty"`
}

// WorkflowService 工作流服务接口
type WorkflowService interface {
	ExecuteWorkflow(ctx context.Context, steps []WorkflowStep) (*WorkflowExecution, error)
	GetWorkflowStatus(ctx context.Context, workflowID string) (*WorkflowExecution, error)
}

// workflowService 工作流服务实现
type workflowService struct {
	log    *log.Helper
	db     *gorm.DB
	active map[string]*WorkflowExecution
}

// NewWorkflowService 创建工作流服务实例
func NewWorkflowService(db *gorm.DB, logger log.Logger) WorkflowService {
	return &workflowService{
		log:    log.NewHelper(logger),
		db:     db,
		active: make(map[string]*WorkflowExecution),
	}
}

// ExecuteWorkflow 执行工作流
func (s *workflowService) ExecuteWorkflow(ctx context.Context, steps []WorkflowStep) (*WorkflowExecution, error) {
	workflowID := fmt.Sprintf("workflow_%s", time.Now().Format("20060102_150405"))
	
	execution := &WorkflowExecution{
		ID:        workflowID,
		Status:    WorkflowRunning,
		Steps:     make([]string, 0, len(steps)),
		Results:   make([]interface{}, 0, len(steps)),
		StartedAt: time.Now(),
	}
	
	// 记录工作流执行
	s.active[workflowID] = execution
	
	s.log.Infof("开始执行工作流: %s", workflowID)
	
	var currentData interface{}
	
	// 按顺序执行每个步骤
	for i, step := range steps {
		execution.Steps = append(execution.Steps, step.GetName())
		
		s.log.Infof("执行步骤 %d/%d: %s", i+1, len(steps), step.GetName())
		
		result, err := step.Execute(ctx, currentData)
		if err != nil {
			execution.Status = WorkflowFailed
			execution.Error = err.Error()
			now := time.Now()
			execution.FinishedAt = &now
			
			s.log.Errorf("工作流执行失败，步骤 %s 出错: %v", step.GetName(), err)
			return execution, err
		}
		
		execution.Results = append(execution.Results, result)
		currentData = result
	}
	
	// 标记工作流完成
	execution.Status = WorkflowCompleted
	now := time.Now()
	execution.FinishedAt = &now
	
	delete(s.active, workflowID)
	
	s.log.Infof("工作流执行完成: %s", workflowID)
	return execution, nil
}

// GetWorkflowStatus 获取工作流状态
func (s *workflowService) GetWorkflowStatus(ctx context.Context, workflowID string) (*WorkflowExecution, error) {
	execution, exists := s.active[workflowID]
	if !exists {
		return nil, fmt.Errorf("工作流不存在: %s", workflowID)
	}
	
	return execution, nil
}