package connectors

type Context struct {
	Name            string     `json:"name"`
	Tenant          string     `json:"tenant"`
	Kubernetes      Connector  `json:"kubernetes"`
	Azure           Connector  `json:"azure"`
	Sharepoint      Sharepoint `json:"sharepoint"`
	MongoDBdatabase string     `json:"mongodbdatabase"`
}

func GetContext() (*Context, error) {
	ctx := Context{}
	// kubes, err := KubernetesClusters()
	// for _, kube := range kubes {
	// 	if kube.IsCurrent {
	// 		ctx.Kubernetes = kube
	// 	}
	// }
	// azure, err := AzureSubscriptions()
	// for _, az := range azure {
	// 	if az.IsCurrent {
	// 		ctx.Azure = az
	// 	}
	// }
	// if err != nil {
	// 	return nil, err
	// }

	mateContext, err := GetMateContext()
	if err != nil {
		return nil, err
	}

	// for _, sp := range mateContext.SharePoint {
	// 	if sp.Tenant == mateContext.Current.SharePoint {
	// 		ctx.Sharepoint = sp
	// 	}
	// }

	// for _, mongo := range mateContext.Mongo {
	// 	if mongo.Cluster == mateContext.Current.Mongo {
	// 		ctx.MongoDBdatabase = mongo.Database
	// 	}
	// }
	ctx.Tenant = mateContext.Current.Tenant
	ctx.Name = mateContext.Current.Name

	return &ctx, nil
}
