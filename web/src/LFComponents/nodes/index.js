// svg png 图片资源来自阿里字体库
// [阿里字体库](https://www.iconfont.cn/collections/index?spm=a313x.7781069.1998910419.4)
// svg图标建议使用自己创建的
// import registerStart from './registerStart'
// import registerUser from './registerUser'
// import registerEnd from './registerEnd'
// import registerPush from './registerPush'
// import registerDownload from './registerDownload'
// import registerPolyline from './registerPolyline'
// import registerTask from './registerTask'
// // import registerConnect from './registerConnect'
// import registerHttp from './registerHttp'
// // import registerApplication from './registerApplication'
import BaseNode from './BaseNode'

function registerEnd(lf){BaseNode(lf, 'end', '/src/assets/icons/end.svg')} 
function registerStart(lf){BaseNode(lf, 'start', '/src/assets/icons/start.svg')}
function registerWeb(lf){BaseNode(lf, 'web', '/src/assets/icons/web2.svg')}
function registerApplication(lf){BaseNode(lf, 'application', '/src/assets/icons/application.svg')}
// jenkins节点
function registerJenkins(lf){BaseNode(lf, 'jenkins', '/src/assets/icons/jenkins.svg')}
export { registerStart, registerEnd,registerApplication,registerWeb,registerJenkins }