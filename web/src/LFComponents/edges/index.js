import {  BezierEdge, BezierEdgeModel } from '@logicflow/core'
import {getShapeStyleFunction} from '../getShapeStyleUtil'
    class myCurvedEdge  extends BezierEdge {
    }
    class myCurvedEdgeModel  extends BezierEdgeModel {
   
  constructor (data, graphModel) {
    super(data, graphModel)
    this.strokeWidth = 1
  }
  getTextStyle () {
    const style = super.getTextStyle()
    return getTextStyleFunction(style, this.properties)
  }
  getEdgeAnimationStyle()
  {
    const style = super.getEdgeAnimationStyle();
    style.stroke = '#FFA940';
    // style.animationName = 'flow'; // 确保和CSS动画名称一致
    style.animationDuration = '10s';
     return style;
  }

  getEdgeStyle () {
    const attributes = super.getEdgeStyle()
    const properties = this.properties;
    const style = getShapeStyleFunction(attributes, properties)
    return { ...style, fill: 'none' }
  }
    }

    export default {
      type: 'myCurvedEdge',
      view: myCurvedEdge,
      model: myCurvedEdgeModel, 
    }