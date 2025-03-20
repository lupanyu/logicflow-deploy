import LogicFlow from '@logicflow/core'
import { Menu, Snapshot, MiniMap, DndPanel } from '@logicflow/extension'
import { nodeList } from '../components/config'
import edges from '../LFComponents/edges/index.js'
import { RegisterNodes} from '../LFComponents/nodes/'

// 公共配置
const defaultOptions = {
  width: 1000,
  height: 600,
  background: {
    backgroundColor: '#f7f9ff',
  },
  grid: {
    size: 10,
    visible: true
  },
  edgeTextDraggable: true,
  hoverOutline: false,
  plugins: [DndPanel, Menu, MiniMap, Snapshot]
}

export function useLogicFlow(containerRef, options = {}) {
  const lf = ref(null)
  
  // 初始化逻辑
  const init = () => {
    lf.value = new LogicFlow({
      ...defaultOptions,
      ...options,
       container: containerRef.value
    })
    
    // 注册公共组件
    RegisterNodes(lf.value)
     lf.value.register(edges)
    lf.value.setDefaultEdgeType('myCurvedEdge')
    return lf.value
  }

 

  return {
    lf,
    init,
   }
}
