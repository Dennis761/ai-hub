import { encrypt } from '../../utils/cryptoUtils.js';

class ApiKeyWriteManager {
  constructor({ apiKeyReadRepository, apiKeyWriteRepository, apiKeyCountRepository }) {
    this.apiKeyReadRepository = apiKeyReadRepository;
    this.apiKeyWriteRepository = apiKeyWriteRepository;
    this.apiKeyCountRepository = apiKeyCountRepository;
  }

  async create(data) {
    const sameKeyNameKeys = await this.apiKeyReadRepository.findMany({ keyName: data.keyName });
  
    const usedByAnotherUser = sameKeyNameKeys.find(
      key => key.ownerId.toString() !== data.ownerId.toString()
    );
  
    if (usedByAnotherUser) {
      throw new Error(`The key name "${data.keyName}" is already used by another user`);
    }
  
    const ownKeys = sameKeyNameKeys.filter(
      key => key.ownerId.toString() === data.ownerId.toString()
    );
  
    if (ownKeys.length >= 3) {
      throw new Error(`Maximum of 3 keys with name "${data.keyName}" allowed`);
    }
  
    const exists = ownKeys.find(key => key.usageEnv === data.usageEnv);
    if (exists) {
      throw new Error(`A key with name "${data.keyName}" already exists in "${data.usageEnv}" environment`);
    }
  
    const encrypted = encrypt(data.keyValue);
    return this.apiKeyWriteRepository.create({ ...data, keyValue: encrypted });
  }
  

  async update(id, updates, adminId) {
    const key = await this.apiKeyReadRepository.findById(id);
    if (!key) throw new Error('API key not found');
    if (key.ownerId.toString() !== adminId.toString()) throw new Error('Forbidden');
  
    const isRenamingKeyName = updates.keyName && updates.keyName !== key.keyName;
    const isChangingEnv = updates.usageEnv && updates.usageEnv !== key.usageEnv;
  
    if (isRenamingKeyName || isChangingEnv) {
      const sameKeyNameKeys = await this.apiKeyReadRepository.findMany({
        keyName: updates.keyName || key.keyName
      });
  
      const usedByAnotherUser = sameKeyNameKeys.find(
        k => k.ownerId.toString() !== adminId.toString()
      );
      if (usedByAnotherUser) {
        throw new Error(`The key name "${updates.keyName || key.keyName}" is already used by another user`);
      }
  
      const ownKeys = sameKeyNameKeys.filter(
        k => k.ownerId.toString() === adminId.toString() && k._id.toString() !== id
      );
  
      if (ownKeys.length >= 3) {
        throw new Error(`Maximum of 3 keys with name "${updates.keyName || key.keyName}" allowed`);
      }
  
      const existsSameEnv = ownKeys.find(k =>
        (updates.usageEnv || key.usageEnv) === k.usageEnv
      );
  
      if (existsSameEnv) {
        throw new Error(`A key with name "${updates.keyName || key.keyName}" already exists in "${updates.usageEnv || key.usageEnv}" environment`);
      }
    }
  
    if (updates.keyValue) {
      updates.keyValue = encrypt(updates.keyValue);
    }
  
    return this.apiKeyWriteRepository.update(id, updates);
  }
  

  async delete(id, adminId) {
    const key = await this.apiKeyReadRepository.findById(id);
    if (!key) throw new Error('API key not found');
    if (key.ownerId.toString() !== adminId.toString()) throw new Error('Forbidden');

    return this.apiKeyWriteRepository.delete(id);
  }
}

export default ApiKeyWriteManager;
