package entity

import "time"

type Project struct {
	Id                    int64
	Name                  string
	Environments          []*Environment
	TeamId                int64
	TeamSlug              string
	PrEnvironmentsEnabled *bool
	ClusterId             *int64
	CreatedAt             time.Time
}

type ProjectCluster struct {
	ClusterId   int64
	ClusterName string
}
