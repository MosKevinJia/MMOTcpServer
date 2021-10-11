// 地图模块
package Game

type MAP struct {
	Width  int
	Height int
	Grid   [][]int
}

// 初始化
func (m *MAP) Init() {

	m.Grid = make([][]int, m.Width)
	for i := range m.Grid {
		m.Grid[i] = make([]int, m.Height)
	}

}

// 获得地图x,y的内容
func (m *MAP) Get(x int, y int) int {
	return m.Grid[x][y]
}
