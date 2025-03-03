export default function registerStart (lf) {
  lf.register('start', ({ CircleNode, CircleNodeModel, h }) => {
    class StartNode extends CircleNode {
      getLabelShape () {
        const {model} = this.props
        const {
          x,
          y
        } = model
        return h(
          'svg',
          {
            fill: '#000000',
            fontSize: 12,
            x: x - 12,
            y: y + 4,
            width: 50,
            height: 25,
            viewBox: '0 0 1024 1024'

          },
          'Start-Node'
        )
      }
      getShape () {
        const {model} = this.props
        const {width, height, x, y, points} = model
        const style = model.getNodeStyle();
        const transform = `matrix(1 0 0 1 ${x - width / 2} ${y - height / 2})`

        return h('g', {transform}, [      
          // 图标部分（关键修改）
          h('svg', {
             x: x - width / 2,
            y: y - height / 2,
            width,
            height,
            rx: points,
            ry: points,
            viewBox: '0 0 1024 1024',
          }), [
            h('image', {
              href: '/src/assets/icons/start.svg',
              fill: style.stroke,
              width,
              height,
            })],
          
          this.getLabelShape()
        ]);
      }
    }

    class StartModel extends CircleNodeModel {
      // 自定义节点形状属性
      initNodeData(data) {
        data.text = {
          value: (data.text && data.text.value) || '',
          x: data.x,
          y: data.y + 35,
          dragable: false,
          editable: true
        }
        super.initNodeData(data)
        this.r = 20
      }
      // 自定义节点样式属性
      getNodeStyle() {
        const style = super.getNodeStyle()
        return style
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
      getConnectedTargetRules () {
        const rules = super.getConnectedTargetRules()
        const notAsTarget = {
          message: '起始节点不能作为连线的终点',
          validate: () => false
        }
        rules.push(notAsTarget)
        return rules
      }
    }
    return {
      view: StartNode,
      model: StartModel
    }
  })
}
