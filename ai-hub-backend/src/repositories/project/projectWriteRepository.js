import ProjectModel from '../../models/projectModel.js';
import TaskModel from '../../models/taskModel.js';
import PromptModel from '../../models/promptModel.js';

export default class ProjectWriteRepository {
  async create(data) {
    const project = new ProjectModel(data);
    return await project.save();
  }

  async update(id, updates) {
    return await ProjectModel.findByIdAndUpdate(id, updates, { new: true });
  }

  async delete(id) {
    const tasks = await TaskModel.find({ projectId: id });
    const taskIds = tasks.map(t => t._id);

    await PromptModel.deleteMany({ taskId: { $in: taskIds } });
    await TaskModel.deleteMany({ projectId: id });

    return await ProjectModel.findByIdAndDelete(id);
  }
}
