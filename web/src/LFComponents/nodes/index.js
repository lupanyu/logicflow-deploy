
import BaseNode from './BaseNode'
import StartNode from './StartNode'
import RectNode from './RectNode'  // 圆角矩形
import CircleNode from './CircleNode' // 圆形

export const ExtendNode=(tpstring,node)=>{
    return {
          type: tpstring,
          view: node.view,
          model: node.model
    }
}

export const RegisterNodes = (lf) => {
 
    lf.register(ExtendNode('web',RectNode))    
    lf.register(ExtendNode('start',CircleNode))    
    lf.register(ExtendNode('end',CircleNode))
    lf.register(ExtendNode('java',RectNode))
    lf.register(ExtendNode('shell',RectNode))
    lf.register(ExtendNode('jenkins',RectNode))

}

function registerEnd(lf){BaseNode(lf, 'end', '/src/assets/icons/end.svg')} 
// function registerStart(lf){BaseNode(lf, 'start', '/src/assets/icons/start.svg')}
function registerStart(lf){BaseNode(lf,'start', '/src/assets/icons/start.svg')}
function registerWeb(lf){BaseNode(lf, 'web', '/src/assets/icons/web2.svg')}
function registerJava(lf){BaseNode(lf, 'java', '/src/assets/icons/application.svg')}
// jenkins节点
function registerJenkins(lf){BaseNode(lf, 'jenkins', '/src/assets/icons/jenkins.svg')}
// shell 节点
function registerShell(lf){BaseNode(lf, 'shell', '/src/assets/icons/shell.svg')}
export { registerStart, registerEnd,registerJava,registerWeb,registerJenkins,registerShell,StartNode}