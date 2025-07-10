import ApiKeyModel from '../../models/apiKeyModel.js';

class ApiKeyWriteRepository {
  async create(data) {
    const key = new ApiKeyModel(data);
    return await key.save();
  }

  async update(id, updates) {
    return await ApiKeyModel.findByIdAndUpdate(id, { $set: updates }, { new: true, runValidators: true });
  }

  async delete(id) {
    return await ApiKeyModel.findByIdAndDelete(id);
  }
}

export default ApiKeyWriteRepository;
