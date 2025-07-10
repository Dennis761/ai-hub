import TaskModel from '../../models/taskModel.js';

class TaskReadRepository {
  findAll(filter = {}) {
    return TaskModel.find(filter).sort({ createdAt: -1 });
  }

  findById(id) {
    return TaskModel.findById(id).populate('prompts');
  }
}

export default TaskReadRepository;
