import PromptModel from '../../models/promptModel.js';

class PromptReadRepository {
  findById(id) {
    return PromptModel.findById(id);
  }

  findByTaskId(taskId) {
    return PromptModel.find({ taskId }).sort({ executionOrder: 1 });
  }
}

export default PromptReadRepository;
