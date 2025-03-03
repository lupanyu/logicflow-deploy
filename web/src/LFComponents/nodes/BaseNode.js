export default function createBaseNode(lf, nodeType, iconPath) {
  return lf.register(nodeType, ({ PolygonNode, PolygonNodeModel, h }) => {
    class Node extends PolygonNode {
      constructor(data) {
        super(data); // 确保父类的构造函数被调用
      }
       getShape() {
        const { model } = this.props;
        const { width, height, points, x, y } = model;
        const style = model.getNodeStyle();
         const transform = `matrix(1 0 0 1 ${x - width / 2} ${y - height / 2})`;
        // 获取文字内容（默认取 model.text.value，如果为空则为空字符串）
        const textContent = model.text && model.text.value ? model.text.value : '';
        // 定义图标和文字的偏移和位置（均为局部坐标）
        const imageX = 0;
        const imageY = 0;
        // 例如：将文字放在图标正下方，文字水平居中
        const textX = width / 2;
        const textY = height + 10; // 根据需要微调偏移值
      
        return h('g', { transform }, [
          // svg 画布从局部 (0,0) 开始绘制
          h('svg', {
            x: 0,
            y: 0,
            width,
            height,
            rx: points,
            ry: points,
            viewBox: "0 0 1028 1024",
          }),
          // 绘制图标图片，使用局部坐标
          h('image', {
            href: iconPath,
            fill: style.stroke,
            x: imageX,
            y: imageY,
            width,
            height,
          }),
          // 绘制文字，坐标设置为局部坐标
          h('text', {
            x: textX,
            y: textY,
            fill: '#000',
            fontSize: 12,
            textAnchor: 'middle',
          }, textContent)
        ]);
      }
      
    }
    class Model extends PolygonNodeModel {
    
      constructor(data, graphModel) {
            super(data, graphModel);
            this.nodeType = nodeType;

        }
        initNodeData(data) {
            super.initNodeData(data);
            // 完全禁用父类文字显示
            this.text.visible = false;
            this.points = [
            [30, 0],
            [60, 30],
            [30, 60],
            [0, 30],
            ];
        }
        // 自定义节点样式属性
        getNodeStyle() {
          const style = super.getNodeStyle()
          return style
        }
        getTextStyle() {
         
          return {
            ...super.getTextStyle(),
            fontSize: 0,
            fill: '#000', 
          }
        }
        
        getAnchorStyle() {
        // 自定义锚点样式
        const HOVER_RADIUS = 8;
        const HOVER_FILL_COLOR = "rgb(24,55,255)";
        const HOVER_STROKE_COLOR = "rgb(24, 125, 255)";
        const style = super.getAnchorStyle();
        style.hover.r = HOVER_RADIUS;
        style.hover.fill = HOVER_FILL_COLOR;
        style.hover.stroke = HOVER_STROKE_COLOR;
        return style;
        }
        // 自定义节点outline
        getOutlineStyle() {
        const style = super.getOutlineStyle();
        style.stroke = '#88f'
        return style
        }
        getConnectedTargetRules() {
             const rules = super.getConnectedTargetRules();
            if (this.nodeType === 'start') {
                const notAsTarget = {
                    message: '起始节点不能作为连线的终点',
                    validate: () => false
                };
                rules.push(notAsTarget);
            } 
            return rules;
            }
        
        getConnectedSourceRules () {
            const rules = super.getConnectedSourceRules()
            if (this.nodeType ==='end') {
                const rules = super.getConnectedSourceRules()
                const notAsTarget = {
                message: '终止节点不能作为连线的起点',
                validate: () => false
                }
            rules.push(notAsTarget)
            }
            return rules
          }
    }

    return {
      type: nodeType,
      view: Node,
      model: Model,
    };
  });
}
