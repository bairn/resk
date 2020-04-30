package infra

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"reflect"
	"sort"
)

const (
	KeyProps = "_conf"
)

type StarterContext map[string]interface{}

func (s StarterContext) Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置还没有被初始化")
	}
	return p.(kvs.ConfigSource)
}

func (s StarterContext) SetProps(conf kvs.ConfigSource) {
	s[KeyProps] = conf
}

type Starter interface {
	Init(StarterContext)
	Setup(StarterContext)
	Start(StarterContext)
	Stop(StarterContext)
	StartBlocking() bool
	PriorityGroup() PriorityGroup
	Priority() int
}


type starterRegister struct {
	nonBlockingStarters []Starter
	blockingStarters []Starter
}

func (r *starterRegister) AllStarters() []Starter {
	starters := make([]Starter, 0)
	starters = append(starters, r.nonBlockingStarters...)
	starters = append(starters, r.blockingStarters...)
	return starters
}

func (r *starterRegister) Register(starter Starter) {
	if starter.StartBlocking() {
		r.blockingStarters = append(r.blockingStarters, starter)
	} else {
		r.nonBlockingStarters = append(r.nonBlockingStarters, starter)
	}

	typ := reflect.TypeOf(starter)
	log.Infof("Register starter:%s", typ.String())
}


var StarterRegister *starterRegister = new(starterRegister)

type Starters []Starter
func (s Starters) Len() int      { return len(s) }
func (s Starters) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Starters) Less(i, j int) bool {
	return s[i].PriorityGroup() > s[j].PriorityGroup() && s[i].Priority() > s[j].Priority()
}

func Register(starter Starter) {
	StarterRegister.Register(starter)
}

func SortStarters() {
	sort.Sort(Starters(StarterRegister.AllStarters()))
}

func GetStarters() []Starter {
	return StarterRegister.AllStarters()
}

type PriorityGroup int

const (
	SystemGroup         PriorityGroup = 30
	BasicResourcesGroup PriorityGroup = 20
	AppGroup            PriorityGroup = 10

	INT_MAX          = int(^uint(0) >> 1)
	DEFAULT_PRIORITY = 10000
)

type BaseStarter struct {
}

func (s *BaseStarter) Init(StarterContext) {}
func (s *BaseStarter) Setup(StarterContext) {}
func (s *BaseStarter) Start(StarterContext) {}
func (s *BaseStarter) Stop(StarterContext) {}
func (s *BaseStarter) StartBlocking() bool {return false}
func (s *BaseStarter) PriorityGroup() PriorityGroup { return  BasicResourcesGroup}
func (s *BaseStarter) Priority() int { return DEFAULT_PRIORITY }