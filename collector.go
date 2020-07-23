package exporter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

const (
	slabinfo             = "/proc/slabinfo"
	namespace            = "slab"
	exporter             = "slab_exporter"
	eventConfigNamespace = "events"
)

var (
	space       = regexp.MustCompile(`\s+`)
	slabVer     = regexp.MustCompile(`slabinfo -`)
	slabHeader  = regexp.MustCompile(`# name`)
	slabMetrics = []string{
		"active_objs",
		"num_objs",
		"objsize",
		"objperslab",
		"pagesperslab",
		"limit",
		"batchcount",
		"sharedfactor",
		"active_slabs",
		"num_slabs",
		"sharedavail",
	}
)

type slabInfo struct {
	name         string
	objActive    int64
	objNum       int64
	objSize      int64
	objPerSlab   int64
	pagesPerSlab int64
	// tunables
	limit        int64
	batch        int64
	sharedFactor int64
	slabActive   int64
	slabNum      int64
	sharedAvail  int64
}

func (s *slabInfo) metrics() []prometheus.Metric {
	return []prometheus.Metric{
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"active",
					"objs",
				),
				fmt.Sprintf(
					"slab %s",
					"active_objs",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.objActive),
			s.name,
		),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"objs",
					"num",
				),
				fmt.Sprintf(
					"slab %s",
					"num_objs",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.objNum),
			s.name,
		),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"obj", "size",
				),
				fmt.Sprintf(
					"slab %s",
					"objsize",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.objSize),
			s.name,
		),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"obj", "perslab",
				),
				fmt.Sprintf(
					"slab %s",
					"objperslab",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.objPerSlab),
			s.name,
		),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"limit",
					"slabs",
				),
				fmt.Sprintf(
					"slab %s",
					"limit",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.limit),
			s.name,
		),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"batch", "count",
				),
				fmt.Sprintf(
					"slab %s",
					"batchcount",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.batch),
			s.name,
		),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"shared", "factor",
				),
				fmt.Sprintf(
					"slab %s",
					"sharedfactor",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.sharedFactor),
			s.name,
		),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"slabs",
					"active",
				),
				fmt.Sprintf(
					"slab %s",
					"active_slabs",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.slabActive),
			s.name,
		),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"slabs",
					"num",
				),
				fmt.Sprintf(
					"slab %s",
					"num_slabs",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.slabNum),
			s.name,
		),
		prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(
					namespace,
					"shared", "avail",
				),
				fmt.Sprintf(
					"slab %s",
					"sharedavail",
				),
				[]string{"slab"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.sharedAvail),
			s.name,
		),
	}
}

func parseSlab(line string) (*slabInfo, error) {
	// First cleanup whitespace.
	l := space.ReplaceAllString(line, " ")
	s := strings.Split(l, " ")
	if len(s) != 16 {
		return nil, fmt.Errorf("unable to parse: %s", line)
	}
	var err error
	i := &slabInfo{name: formatName(s[0])}
	i.objActive, err = strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		return nil, err
	}
	i.objNum, err = strconv.ParseInt(s[2], 10, 64)
	if err != nil {
		return nil, err
	}
	i.objSize, err = strconv.ParseInt(s[3], 10, 64)
	if err != nil {
		return nil, err
	}
	i.objPerSlab, err = strconv.ParseInt(s[4], 10, 64)
	if err != nil {
		return nil, err
	}
	i.pagesPerSlab, err = strconv.ParseInt(s[5], 10, 64)
	if err != nil {
		return nil, err
	}
	i.limit, err = strconv.ParseInt(s[8], 10, 64)
	if err != nil {
		return nil, err
	}
	i.batch, err = strconv.ParseInt(s[9], 10, 64)
	if err != nil {
		return nil, err
	}
	i.sharedFactor, err = strconv.ParseInt(s[10], 10, 64)
	if err != nil {
		return nil, err
	}
	i.slabActive, err = strconv.ParseInt(s[13], 10, 64)
	if err != nil {
		return nil, err
	}
	i.slabNum, err = strconv.ParseInt(s[14], 10, 64)
	if err != nil {
		return nil, err
	}
	i.sharedAvail, err = strconv.ParseInt(s[15], 10, 64)
	if err != nil {
		return nil, err
	}
	return i, nil
}

type slabCollector struct {
	r *regexp.Regexp
}

// SlabCollector is a prometheus collector.
type SlabCollector interface {
	prometheus.Collector
}

// NewSlabCollector returns a new slab collector.
func NewSlabCollector(config *viper.Viper) (SlabCollector, error) {
	var (
		r   *regexp.Regexp
		err error
	)
	rStr := config.GetString("regex")
	if rStr != "" {
		r, err = regexp.Compile(rStr)
		if err != nil {
			return nil, err
		}
	}
	return &slabCollector{
		r: r,
	}, nil
}

// Describe implements the prometheus.Collector interface.
func (c *slabCollector) Describe(ch chan<- *prometheus.Desc) {
	f, err := os.Open(slabinfo)
	if err != nil {
		log.Println(err)
		return
	}

	slabInfos := []*slabInfo{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !c.shouldParse(line) {
			continue
		}
		s, err := parseSlab(line)
		if err != nil {
			log.Println(err)
			continue
		}
		slabInfos = append(slabInfos, s)
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	sort.Slice(slabInfos, func(i, j int) bool {
		return slabInfos[i].name > slabInfos[j].name
	})
	for _, m := range slabMetrics {
		ch <- prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				m,
				"slab",
			),
			fmt.Sprintf(
				"slab %s",
				m,
			),
			[]string{"slab"},
			nil,
		)
	}
}

// Collect implements prometheus.Collector interface.
func (c *slabCollector) Collect(ch chan<- prometheus.Metric) {
	f, err := os.Open(slabinfo)
	if err != nil {
		log.Println(err)
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !c.shouldParse(line) {
			continue
		}
		s, err := parseSlab(line)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, m := range s.metrics() {
			ch <- m
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func (c *slabCollector) shouldParse(line string) bool {
	if slabVer.MatchString(line) {
		return false
	}
	if slabHeader.MatchString(line) {
		return false
	}
	if c.r != nil {
		if c.r.MatchString(line) {
			return true
		}
		return false
	}
	return true
}

func formatName(s string) string {
	return strings.Replace(strings.Replace(strings.Replace(
		strings.Replace(s, ":", "_", -1),
		".",
		"_",
		-1,
	),
		"-",
		"_",
		-1,
	),
		"/",
		"_",
		-1,
	)
}
