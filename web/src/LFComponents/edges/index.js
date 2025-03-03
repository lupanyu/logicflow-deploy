import { CurvedEdge, CurvedEdgeModel } from '@logicflow/extension'

    class myCurvedEdge  extends CurvedEdge {
    }
    class myCurvedEdgeModel  extends CurvedEdgeModel {
      getEdgeStyle() {
        const style = super.getEdgeStyle();
        //  正确获取color属性
        const color = this.properties.color || 'default';
        const colorMap = {
          default: '#999999',
          green: '#52c41a',
          red: '#ff4d4f',
          blue: '#1890ff'
        };
        style.stroke = colorMap[color];
        style.strokeWidth = 2;
        return style;
      }
      getEdgeAnimationStyle() {
        const style = super.getEdgeAnimationStyle();
        style.stroke = 'green';
        style.animationDuration = '30s';
        
        // style.animationDirection = 'reverse';
        return style;
      }

      initEdgeData(data) {
        super.initEdgeData(data);
        this.radius = 20;

        // 确保初始化时设置默认属性
        if (!data.properties?.color) {
          this.setProperties({ color: 'default' });
        }
      }

      // 监听属性变化
      propertiesChanged(val) {
        super.propertiesChanged(val);
        this.updateStroke();
      }

      updateStroke() {
        const newStyle = this.getEdgeStyle();
        this.stroke = newStyle.stroke;
        this.strokeWidth = newStyle.strokeWidth;
      }
    }

    export default {
      type: 'myCurvedEdge',
      view: myCurvedEdge,
      model: myCurvedEdgeModel, 
    }