package clusters

//Returns the oldest N Items
func (l *ClusterList) OldestN(count int) *ClusterList {
	oldestlist := NewClusterList()
	if len(l.items) > count {
		for i := 0; i < count; i++ {
			var oldest *Cluster
			oldestassigned := false
			for _, item := range l.items {
				if !oldestlist.ClusterInList(item) {
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
			oldestlist.items = append(oldestlist.items, oldest)
		}
	} else {
		oldestlist =  l
	}
	return oldestlist
}

//ClusterInList checks if specified cluster is in specified list
func (l *ClusterList) ClusterInList(c *Cluster) bool {
	for _, item := range l.items {
		if Equal(c, item) {
			return true
		}
	}
	return false
}


//Extracts list of Items in a specific state
func (l *ClusterList) ClustersInStateIn(state string) *ClusterList {
	oc := NewClusterList()
	for _, c := range l.items {
		if c.State == state {
			oc.items = append(oc.items, c)
		}
	}
	return oc
}


//List iterates over all the clusters
func (l *ClusterList) List(filter func(c *Cluster))  {
	for _, item := range l.items {
		filter(item)
	}
}

//Appends a cluster to the list
func (l *ClusterList) Append(c *Cluster)  {
	l.items = append(l.items, c)
	return
}

//Len gets the length of the clusters
func (l *ClusterList) Len() int {
	return len(l.items)
}

//ItemAt gets cluster of specified index
func (l *ClusterList) ItemAt(index int) *Cluster {
	return l.items[index]
}