import CircleNewModel from './CircleNode'

class StartNodeModel extends CircleNewModel.model {
}

export default {
    type: 'end',
    view: CircleNewModel.view,
    model: StartNodeModel 
}