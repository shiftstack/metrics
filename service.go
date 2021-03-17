package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type versionInformation struct {
	ID         string
	MinVersion string
	Version    string
}

func (v versionInformation) Labels() map[string]string {
	return map[string]string{
		"id":          v.ID,
		"min_version": v.MinVersion,
		"version":     v.Version,
	}
}

type service struct {
	*prometheus.GaugeVec
	Versions func() ([]versionInformation, error)
}

func (s service) Refresh() error {
	versions, err := s.Versions()
	if err != nil {
		return err
	}

	s.Reset()
	for _, v := range versions {
		s.With(v.Labels()).Set(1)
	}

	return nil
}

var versionLabels = []string{"id", "min_version", "version"}

func NewService(name, help string, versionsFunc func() ([]versionInformation, error)) service {
	return service{
		GaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, versionLabels),

		Versions: versionsFunc,
	}
}
