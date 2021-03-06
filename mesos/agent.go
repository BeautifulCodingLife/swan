package mesos

import (
	"encoding/json"
	"sync"

	"github.com/Dataman-Cloud/swan/mesosproto"
)

type Agent struct {
	id       string
	hostname string
	attrs    []*mesosproto.Attribute

	sync.RWMutex
	offers map[string]*Offer
	tasks  map[string]*Task
}

func newAgent(id, hostname string, attrs []*mesosproto.Attribute) *Agent {
	return &Agent{
		id:       id,
		hostname: hostname,
		attrs:    attrs,
		offers:   make(map[string]*Offer),
		tasks:    make(map[string]*Task),
	}
}

func (s *Agent) MarshalJSON() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	m := map[string]interface{}{
		"id":         s.id,
		"hostname":   s.hostname,
		"attributes": s.attrs,
		"offers":     s.offers,
	}
	return json.Marshal(m)
}

func (s *Agent) addOffer(offer *Offer) {
	s.Lock()
	s.offers[offer.GetId()] = offer
	s.Unlock()
}

func (s *Agent) removeOffer(offerID string) bool {
	s.Lock()
	defer s.Unlock()

	found := false
	if _, found = s.offers[offerID]; found {
		delete(s.offers, offerID)
	}
	return found
}

func (s *Agent) empty() bool {
	s.RLock()
	defer s.RUnlock()

	return len(s.offers) == 0 && len(s.tasks) == 0
}

func (s *Agent) getOffer(offerId string) *Offer {
	s.RLock()
	defer s.RUnlock()

	offer, ok := s.offers[offerId]
	if !ok {
		return nil
	}

	return offer
}

func (s *Agent) getOffers() []*Offer {
	s.RLock()
	defer s.RUnlock()

	offers := make([]*Offer, 0)
	for _, offer := range s.offers {
		offers = append(offers, offer)
	}

	return offers
}

func (s *Agent) offer() *Offer {
	s.RLock()
	defer s.RUnlock()

	for _, offer := range s.offers {
		return offer
	}

	return nil
}

func (s *Agent) addTask(t *Task) {
	s.Lock()
	defer s.Unlock()

	s.tasks[t.TaskId.GetValue()] = t
}

func (s *Agent) getTask(taskId string) *Task {
	s.RLock()
	defer s.RUnlock()

	task, ok := s.tasks[taskId]
	if !ok {
		return nil
	}

	return task
}

func (s *Agent) removeTask(taskId string) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.tasks[taskId]
	if !ok {
		return
	}

	delete(s.tasks, taskId)
}

func (s *Agent) getTasks() []*Task {
	s.RLock()
	defer s.RUnlock()

	tasks := make([]*Task, 0)
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (s *Agent) Resources() (cpus, mem, disk float64, ports []uint64) {
	for _, offer := range s.getOffers() {
		cpus += offer.GetCpus()
		mem += offer.GetMem()
		disk += offer.GetDisk()
		ports = append(ports, offer.GetPorts()...)
	}

	return
}

func (s *Agent) Attributes() map[string]string {
	attrs := make(map[string]string)

	for _, offer := range s.getOffers() {
		for k, v := range offer.GetAttrs() {
			attrs[k] = v
		}
	}

	return attrs
}
