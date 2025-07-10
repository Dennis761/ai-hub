import ApiKeyModel from '../../models/apiKeyModel.js';

class ApiKeyReadRepository {
  async findAll(filter = {}) {
    return await ApiKeyModel.find(filter).sort({ createdAt: -1 });
  }

  async findById(id) {
    return await ApiKeyModel.findById(id);
  }

  async findOne(filter) {
    return await ApiKeyModel.findOne(filter);
  }

  async findMany(filter) {
    return await ApiKeyModel.find(filter);
  }
}

export default ApiKeyReadRepository;
