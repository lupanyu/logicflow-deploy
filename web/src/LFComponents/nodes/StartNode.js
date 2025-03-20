import CircleNewModel from './CircleNode'

class StartNodeModel extends CircleNewModel.model {
}

export default {
    type: 'start',
    view: CircleNewModel.view,
    model: StartNodeModel 
}