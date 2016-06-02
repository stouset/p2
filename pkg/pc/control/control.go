/*
Package control wraps a pcstore with convenience functions suitable for operational tasks.
*/
package control

import (
	"github.com/square/p2/pkg/kp"
	"github.com/square/p2/pkg/kp/pcstore"
	"github.com/square/p2/pkg/pc/fields"
	"github.com/square/p2/pkg/types"
	"github.com/square/p2/pkg/util"
	"k8s.io/kubernetes/pkg/labels"
)

type PodCluster struct {
	pcStore pcstore.Store

	az       fields.AvailabilityZone
	cn       fields.ClusterName
	podID    types.PodID
	selector labels.Selector
	session  kp.Session
}

func NewPodCluster(
	az fields.AvailabilityZone,
	cn fields.ClusterName,
	podID types.PodID,
	pcstore pcstore.Store,
	selector labels.Selector,
	session kp.Session,
) *PodCluster {

	pc := &PodCluster{}
	pc.az = az
	pc.cn = cn
	pc.podID = podID
	pc.pcStore = pcstore
	pc.selector = selector
	pc.session = session

	return pc
}

func (pccontrol *PodCluster) All() ([]fields.PodCluster, error) {
	return pccontrol.pcStore.FindWhereLabeled(pccontrol.podID, pccontrol.az, pccontrol.cn)
}

// Best effort delete of the list of podClusterID will not halt on error
func (pccontrol *PodCluster) Delete() (errors []error) {
	podClusterIDs, err := pccontrol.All()
	if err != nil {
		return []error{err}
	}

	for _, pc := range podClusterIDs {
		if err := pccontrol.pcStore.Delete(pc.ID); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (pccontrol *PodCluster) Create(annotations fields.Annotations) (fields.PodCluster, error) {
	return pccontrol.pcStore.Create(pccontrol.podID, pccontrol.az, pccontrol.cn, pccontrol.selector, annotations, pccontrol.session)
}

func (pccontrol *PodCluster) Get() (fields.PodCluster, error) {
	labeledPCs, err := pccontrol.All()
	if err != nil {
		return fields.PodCluster{}, err
	}
	if len(labeledPCs) > 1 {
		return fields.PodCluster{}, util.Errorf("More than one PC matches this podCluster %+v", pccontrol)
	}
	if len(labeledPCs) == 0 {
		return fields.PodCluster{}, util.Errorf("No pod cluster found for this podcluster: %+v", pccontrol)
	}

	return pccontrol.pcStore.Get(labeledPCs[0].ID)
}
