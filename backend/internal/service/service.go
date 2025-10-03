package service
/* *
 * @Author: chengjiang
 * @Date: 2025-10-02 15:21:17
 * @Description: 
**/

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGreeterService, 
	NewUserFeedbackService,
	NewChatService,
	NewWorkflowAPIService,
	NewReportService,
)
