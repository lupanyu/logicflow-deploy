export default function registerHttp(lf) {
  lf.register('application', ({ PolygonNode, PolygonNodeModel, h }) => {
    class Node extends PolygonNode {
        getShape() {
            const {model} = this.props
            const {width, height, x, y, points} = model
            const style = model.getNodeStyle();
            const transform = `matrix(1 0 0 1 ${x - width / 2} ${y - height / 2})`
             return h('g', {transform}, [                     
             h('svg', {
              x: x - width / 2,
              y: y - height / 2,
              width,
              height,
              rx: points,
              ry: points,
              viewBox: "0 0 1028 1024",

            }), [
              h('image', {
                href: '/src/assets/icons/application.svg',
                fill: style.stroke,
                width,
                height,
              })]
            ]);
          }
    }
    class Model extends PolygonNodeModel {
      constructor(data, graphModel) {
        super(data, graphModel);
      }
      initNodeData(data) {
        super.initNodeData(data);
        this.points = [
          [50, 0],
          [100, 50],
          [50, 100],
          [0, 50]
        ];
      }
    }
    return {
      type: 'application',
      view: Node,
      model: Model
    };
  });
} 