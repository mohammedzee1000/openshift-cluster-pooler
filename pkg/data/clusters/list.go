package clusters

//Returns the oldest N Items
func (l ClusterList) OldestN(count int) *ClusterList {
	oldestlist := NewClusterList()
	if len(l.Items) > count {
		for i := 0; i < count; i++ {
			var oldest Cluster
			oldestassigned := false
			for _, item := range l.Items {
				if !oldestlist.ClusterInList(&item) {
					if ! oldestassigned {
						oldest = item
						oldestassigned = true
					} else {
						if item.CreatedOn.After(oldest.CreatedOn) {
							oldest = item
						}
					}
				}
			}
			oldestlist.Items = append(oldestlist.Items, oldest)
		}
	} else {
		oldestlist.Items = l.Items
	}
	return oldestlist
}

//ClusterInList checks if specified cluster is in specified list
func (l ClusterList) ClusterInList(c *Cluster) bool {
	for _, item := range l.Items {
		if Equal(c, &item) {
			return true
		}
	}
	return false
}


//Extracts list of Items in a specific state
func (l ClusterList) ClustersInStateIn(state string) *ClusterList {
	oc := NewClusterList()
	for _, c := range l.Items {
		if c.State == state {
			oc.Items = append(oc.Items, c)
		}
	}
	return oc
}
