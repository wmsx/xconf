package handler

import (
	"context"

	"github.com/micro/go-micro/v2/util/log"
	"github.com/wmsx/xconf/config-srv/dao"
	"github.com/wmsx/xconf/proto/config"
)

func (c *Config) CreateNamespace(ctx context.Context, req *config.NamespaceRequest, rsp *config.NamespaceResponse) error {
	namespace, err := dao.GetDao().CreateNamespace(
		req.GetAppName(),
		req.GetClusterName(),
		req.GetNamespaceName(),
		req.GetFormat(),
		req.GetDescription())
	if err != nil {
		log.Error("[CreateNamespace]", err)
		return err
	}

	rsp.Id = int64(namespace.ID)
	rsp.CreatedAt = namespace.CreatedAt.Unix()
	rsp.UpdatedAt = namespace.UpdatedAt.Unix()
	rsp.AppName = namespace.AppName
	rsp.ClusterName = namespace.ClusterName
	rsp.NamespaceName = namespace.NamespaceName
	rsp.Format = namespace.Format
	rsp.Value = namespace.Value
	rsp.Released = namespace.Released
	rsp.EditValue = namespace.EditValue
	rsp.Description = namespace.Description
	return nil
}

func (c *Config) QueryNamespace(ctx context.Context, req *config.NamespaceRequest, rsp *config.NamespaceResponse) error {
	namespace, err := dao.GetDao().QueryNamespace(req.GetAppName(), req.GetClusterName(), req.GetNamespaceName())
	if err != nil {
		log.Error("[QueryNamespace]", err)
		return err
	}

	rsp.Id = int64(namespace.ID)
	rsp.AppName = namespace.AppName
	rsp.ClusterName = namespace.ClusterName
	rsp.NamespaceName = namespace.NamespaceName
	rsp.Description = namespace.Description
	rsp.Value = namespace.Value
	rsp.Format = namespace.Format
	rsp.EditValue = namespace.EditValue
	rsp.Released = namespace.Released
	rsp.CreatedAt = namespace.CreatedAt.Unix()
	rsp.UpdatedAt = namespace.UpdatedAt.Unix()
	return nil
}

func (c *Config) DeleteNamespace(ctx context.Context, req *config.NamespaceRequest, rsp *config.Response) (err error) {
	err = dao.GetDao().DeleteNamespace(req.GetAppName(), req.GetClusterName(), req.GetNamespaceName())
	if err != nil {
		log.Error("[DeleteNamespace] delete namespace:%s-%s-%s error: %s", req.GetAppName(), req.GetClusterName(), req.GetNamespaceName(), err.Error())
	}
	return
}

func (c *Config) ListNamespaces(ctx context.Context, req *config.ClusterRequest, rsp *config.NamespacesResponse) error {
	namespaces, err := dao.GetDao().ListNamespaces(req.GetAppName(), req.GetClusterName())
	if err != nil {
		log.Error("[ListNamespaces]", err)
		return err
	}

	for _, v := range namespaces {
		rsp.Namespaces = append(rsp.Namespaces, &config.NamespaceResponse{
			Id:            int64(v.ID),
			CreatedAt:     v.CreatedAt.Unix(),
			UpdatedAt:     v.UpdatedAt.Unix(),
			AppName:       v.AppName,
			ClusterName:   v.ClusterName,
			NamespaceName: v.NamespaceName,
			Format:        v.Format,
			Value:         v.Value,
			Released:      v.Released,
			EditValue:     v.EditValue,
			Description:   v.Description,
		})
	}
	return nil
}
