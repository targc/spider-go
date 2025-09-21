package usecase

import (
	"github.com/targc/spider-go/pkg/spider"
	"github.com/google/uuid"
)

type CreateFlowRequest struct {
	TenantID string                 `json:"tenant_id"`
	Name     string                 `json:"name"`
	Meta     map[string]string      `json:"meta,omitempty"`
	Actions  []WorkflowActionInput  `json:"actions"`
	Peers    []PeerInput            `json:"peers"`
}

type WorkflowActionInput struct {
	Key      string                   `json:"key"`
	ActionID string                   `json:"action_id"`
	Config   map[string]string        `json:"config"`
	Mapper   map[string]spider.Mapper `json:"mapper"`
	Meta     map[string]string        `json:"meta,omitempty"`
}

type PeerInput struct {
	ParentKey  string `json:"parent_key"`
	MetaOutput string `json:"meta_output"`
	ChildKey   string `json:"child_key"`
}

type FlowResponse struct {
	FlowID   string `json:"flow_id"`
	FlowName string `json:"flow_name"`
}

func (u *Usecase) CreateFlow(req *CreateFlowRequest) (*FlowResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	flow, err := u.storage.CreateFlow(u.ctx, &spider.CreateFlowRequest{
		ID:       id.String(),
		TenantID: req.TenantID,
		Name:     req.Name,
		Meta:     req.Meta,
	})

	if err != nil {
		return nil, err
	}

	workflowID := flow.ID

	for _, action := range req.Actions {
		_, err = u.storage.AddAction(u.ctx, &spider.AddActionRequest{
			TenantID:   req.TenantID,
			WorkflowID: workflowID,
			Key:        action.Key,
			ActionID:   action.ActionID,
			Config:     action.Config,
			Map:        action.Mapper,
			Meta:       action.Meta,
		})

		if err != nil {
			return nil, err
		}
	}

	for _, peer := range req.Peers {
		err = u.storage.AddDep(
			u.ctx,
			req.TenantID,
			workflowID,
			peer.ParentKey,
			peer.MetaOutput,
			peer.ChildKey,
		)

		if err != nil {
			return nil, err
		}
	}

	return &FlowResponse{
		FlowID:   workflowID,
		FlowName: flow.Name,
	}, nil
}

func (u *Usecase) ListFlows(tenantID string, page, pageSize int) (*spider.FlowListResponse, error) {
	return u.storage.ListFlows(u.ctx, tenantID, page, pageSize)
}

type FlowDetailResponse struct {
	FlowID   string                  `json:"flow_id"`
	FlowName string                  `json:"flow_name"`
	TenantID string                  `json:"tenant_id"`
	Actions  []spider.WorkflowAction `json:"actions"`
}

func (u *Usecase) GetFlow(tenantID, flowID string) (*FlowDetailResponse, error) {
	flow, err := u.storage.GetFlow(u.ctx, tenantID, flowID)
	if err != nil {
		return nil, err
	}

	actions, err := u.storage.GetWorkflowActions(u.ctx, tenantID, flowID)
	if err != nil {
		return nil, err
	}

	return &FlowDetailResponse{
		FlowID:   flowID,
		FlowName: flow.Name,
		TenantID: tenantID,
		Actions:  actions,
	}, nil
}

func (u *Usecase) UpdateFlow(tenantID, flowID string, name string, meta map[string]string) (*spider.Flow, error) {
	req := &spider.UpdateFlowRequest{
		TenantID: tenantID,
		FlowID:   flowID,
		Name:     name,
		Meta:     meta,
	}
	return u.storage.UpdateFlow(u.ctx, req)
}

func (u *Usecase) DeleteFlow(tenantID, flowID string) error {
	return u.storage.DeleteFlow(u.ctx, tenantID, flowID)
}