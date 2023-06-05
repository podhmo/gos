package M



// 順序
type Ordering string

const (
    
    // 降順
    OrderingDesc Ordering = "desc" // default
    
    // 昇順
    OrderingAsc Ordering = "asc"
)



type Season int

const (
    
    SeasonSpring Season = 0 // default
    
    SeasonSummer Season = 1
    
    SeasonAutumn Season = 2
    
    SeasonWinter Season = 3
)
