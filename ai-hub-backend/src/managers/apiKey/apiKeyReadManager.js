import { decrypt } from '../../utils/cryptoUtils.js';

class ApiKeyReadManager {
  constructor(apiKeyReadRepository) {
    this.apiKeyReadRepository = apiKeyReadRepository;
  }

  async getAllKeys(filter = {}) {
    const keys = await this.apiKeyReadRepository.findAll(filter);
    return keys.map(key => ({
      ...key.toObject(),
      keyValue: decrypt(key.keyValue),
    }));
  }

  async getDecryptedKeyById(id) {
    const model = await this.apiKeyReadRepository.findById(id);
    if (!model) throw new Error('Model not found with given ID');
    return {
      ...model.toObject(),
      keyValue: decrypt(model.keyValue),
    };
  }
}

export default ApiKeyReadManager;
