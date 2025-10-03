package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// MCPService MCP服务集成接口
type MCPService interface {
	// QueryDatabase 执行数据库查询
	QueryDatabase(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error)
	// GenerateReport 生成报表
	GenerateReport(ctx context.Context, reportType string, filters map[string]interface{}) ([]byte, error)
	// GetDatabaseSchema 获取数据库结构信息
	GetDatabaseSchema(ctx context.Context) ([]map[string]interface{}, error)
}

// mcpService MCP服务实现
type mcpService struct {
	db  *gorm.DB
	log *log.Helper
}

// NewMCPService 创建MCP服务实例
func NewMCPService(db *gorm.DB, logger log.Logger) MCPService {
	return &mcpService{
		db:  db,
		log: log.NewHelper(logger),
	}
}

// QueryDatabase 执行数据库查询
func (s *mcpService) QueryDatabase(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	s.log.WithContext(ctx).Infof("Executing MCP database query: %s", query)

	var results []map[string]interface{}
	err := s.db.WithContext(ctx).Raw(query, params).Scan(&results).Error
	if err != nil {
		s.log.WithContext(ctx).Errorf("MCP database query failed: %v", err)
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	s.log.WithContext(ctx).Infof("MCP query executed successfully, returned %d rows", len(results))
	return results, nil
}

// GenerateReport 生成报表
func (s *mcpService) GenerateReport(ctx context.Context, reportType string, filters map[string]interface{}) ([]byte, error) {
	s.log.WithContext(ctx).Infof("Generating MCP report of type: %s", reportType)

	switch reportType {
	case "user_feedback_summary":
		return s.generateUserFeedbackReport(ctx, filters)
	case "database_statistics":
		return s.generateDatabaseStatisticsReport(ctx, filters)
	default:
		err := fmt.Errorf("unsupported report type: %s", reportType)
		s.log.WithContext(ctx).Errorf("MCP report generation failed: %v", err)
		return nil, err
	}
}

// generateUserFeedbackReport 生成用户反馈报表
func (s *mcpService) generateUserFeedbackReport(ctx context.Context, filters map[string]interface{}) ([]byte, error) {
	var totalCount int64
	var pendingCount int64
	var avgRating float64

	// 统计总反馈数
	if err := s.db.WithContext(ctx).Table("user_feedbacks").Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count total feedbacks: %w", err)
	}

	// 统计待处理反馈数
	if err := s.db.WithContext(ctx).Table("user_feedbacks").Where("status = ?", "pending").Count(&pendingCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending feedbacks: %w", err)
	}

	// 计算平均评分
	if err := s.db.WithContext(ctx).Table("user_feedbacks").Select("AVG(rating)").Row().Scan(&avgRating); err != nil {
		avgRating = 0
	}

	report := fmt.Sprintf(`
User Feedback Report
Generated at: %s

Summary:
- Total Feedbacks: %d
- Pending Feedbacks: %d
- Average Rating: %.2f

Status Distribution:
- Pending: %d
- Resolved: %d
- Closed: %d
`,
		time.Now().Format(time.RFC3339),
		totalCount,
		pendingCount,
		avgRating,
		pendingCount,
		totalCount-pendingCount,
		0) // 这里可以根据实际状态进一步细化

	return []byte(report), nil
}

// generateDatabaseStatisticsReport 生成数据库统计报表
func (s *mcpService) generateDatabaseStatisticsReport(ctx context.Context, filters map[string]interface{}) ([]byte, error) {
	var tableCount int64
	var totalRecords int64

	// 获取表数量
	if err := s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount).Error; err != nil {
		return nil, fmt.Errorf("failed to get table count: %w", err)
	}

	// 获取总记录数（示例：只统计user_feedbacks表）
	if err := s.db.WithContext(ctx).Table("user_feedbacks").Count(&totalRecords).Error; err != nil {
		totalRecords = 0
	}

	report := fmt.Sprintf(`
Database Statistics Report
Generated at: %s

Database Overview:
- Total Tables: %d
- Total Records (user_feedbacks): %d

Table Information:
- user_feedbacks: %d records
`,
		time.Now().Format(time.RFC3339),
		tableCount,
		totalRecords,
		totalRecords)

	return []byte(report), nil
}

// GetDatabaseSchema 获取数据库结构信息
func (s *mcpService) GetDatabaseSchema(ctx context.Context) ([]map[string]interface{}, error) {
	s.log.WithContext(ctx).Info("Retrieving database schema via MCP")

	var tables []map[string]interface{}
	err := s.db.WithContext(ctx).Raw(`
		SELECT 
			table_name,
			table_type,
			COUNT(column_name) as column_count
		FROM information_schema.tables t
		LEFT JOIN information_schema.columns c ON t.table_name = c.table_name
		WHERE t.table_schema = 'public'
		GROUP BY table_name, table_type
		ORDER BY table_name
	`).Scan(&tables).Error

	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to retrieve database schema: %v", err)
		return nil, fmt.Errorf("schema retrieval failed: %w", err)
	}

	s.log.WithContext(ctx).Infof("Retrieved schema for %d tables", len(tables))
	return tables, nil
}