package clusters

//Returns the oldest N clusters
func OldestN(clusters []Cluster, count int) []Cluster {
	var oldestlist []Cluster
	for i := 0; i < count; i++ {
		var oldest Cluster
		oldestassigned := false
		for _, item := range clusters {
			if !ClusterInList(oldestlist, &item) {
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
		oldestlist = append(oldestlist, oldest)
	}
	return oldestlist
}

func ClusterInList(clusters []Cluster, c *Cluster) bool {
	for _, item := range clusters{
		if Equal(c, &item) {
			return true
		}
	}
	return false
}

//Extracts list of clusters in a specific state
func ClustersInStateIn(inc []Cluster, state string) []Cluster {
	var oc []Cluster
	for _, c := range inc {
		if c.State == state {
			oc = append(oc, c)
		}
	}
	return oc
}