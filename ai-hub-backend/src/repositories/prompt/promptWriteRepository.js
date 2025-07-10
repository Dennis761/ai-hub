import PromptModel from '../../models/promptModel.js';

class PromptWriteRepository {
  create(data) {
    return PromptModel.create(data);
  }

  update(id, updates) {
    return PromptModel.findByIdAndUpdate(id, updates, { new: true });
  }

  delete(id) {
    return PromptModel.findByIdAndDelete(id);
  }

  save(prompt) {
    return prompt.save();
  }
}

export default PromptWriteRepository;
