import TaskModel from '../../models/taskModel.js';
import PromptModel from '../../models/promptModel.js';

class TaskWriteRepository {
  async create(data) {
    return TaskModel.create(data);
  }

  async update(id, updates) {
    return TaskModel.findByIdAndUpdate(id, updates, { new: true });
  }

  async delete(id) {
    const task = await TaskModel.findById(id);
    if (!task) return null;

    await PromptModel.deleteMany({ taskId: task._id });

    await TaskModel.findByIdAndDelete(id);

    return task;
  }
}

export default TaskWriteRepository;
