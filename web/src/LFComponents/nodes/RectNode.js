import { RectResize } from '@logicflow/extension'
import { getShapeStyleFunction, getTextStyleFunction } from '../getShapeStyleUtil'

// 矩形
class RectNewModel extends RectResize.model {

  setToBottom () {
    this.zIndex = 0
  }
  // 圆角
  setAttributes () {
    super.setAttributes()
    this.radius = 20
  }

  getNodeStyle () {
    const style = super.getNodeStyle()
    const properties = this.getProperties()
    return getShapeStyleFunction(style, properties)
  }

  getTextStyle () {
    const style = super.getTextStyle()
    const properties = this.getProperties()
    return getTextStyleFunction(style, properties)
  }
}

export default {
  type: 'pro-rect',
  view: RectResize.view,
  model: RectNewModel
}