package usecase

import "github.com/targc/spider-go/pkg/spider"

func (u *Usecase) DisableAction(tenantID, workflowID, key string) error {
	return u.storage.DisableWorkflowAction(u.ctx, tenantID, workflowID, key)
}

func (u *Usecase) UpdateAction(req *spider.UpdateActionRequest) (*spider.WorkflowAction, error) {
	return u.storage.UpdateAction(u.ctx, req)
}
