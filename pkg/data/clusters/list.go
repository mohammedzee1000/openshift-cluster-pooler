package clusters

func (l ClusterList) Append(cluster Cluster)  {
	l.items = append(l.items, cluster)
}

//Returns the oldest N items
func (l ClusterList) OldestN(count int) *ClusterList {
	oldestlist := NewClusterList()
	for i := 0; i < count; i++ {
		var oldest Cluster
		oldestassigned := false
		for _, item := range l.items {
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
		oldestlist.Append(oldest)
	}
	return oldestlist
}

//ClusterInList checks if specified cluster is in specified list
func (l ClusterList) ClusterInList(c *Cluster) bool {
	for _, item := range l.items {
		if Equal(c, &item) {
			return true
		}
	}
	return false
}


//Extracts list of items in a specific state
func (l ClusterList) ClustersInStateIn(state string) *ClusterList {
	oc := NewClusterList()
	for _, c := range l.items {
		if c.State == state {
			oc.Append(c)
		}
	}
	return oc
}

func (l ClusterList) Rangefunc(f func(c *Cluster))  {
	for _, item := range l.items {
		f(&item)
	}
}

func (l ClusterList) Len() int  {
	return len(l.items)
}

func (l ClusterList) Get(index int) *Cluster {
	return &l.items[index]
}

func (l ClusterList) GetItems() []Cluster  {
	return l.items
}