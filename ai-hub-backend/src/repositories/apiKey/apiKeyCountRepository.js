import ApiKeyModel from '../../models/apiKeyModel.js';

class ApiKeyCountRepository {
  async countKeys(filter) {
    return await ApiKeyModel.countDocuments(filter);
  }
}

export default ApiKeyCountRepository;
