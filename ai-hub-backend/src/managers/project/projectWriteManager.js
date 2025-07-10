import { encrypt } from '../../utils/cryptoUtils.js';

export default class ProjectWriteManager {
  constructor({projectWriteRepository, projectReadRepository}) {
    this.projectWriteRepository = projectWriteRepository;
    this.projectReadRepository = projectReadRepository;
  }

  async create(data) {
    const encryptedKey = encrypt(data.apiKey);
    return await this.projectWriteRepository.create({ ...data, apiKey: encryptedKey });
  }

  async update(id, updates) {
    if (updates.name) {
      const existing = await this.projectReadRepository.findByName(updates.name);
      if (existing && existing._id.toString() !== id) {
        const error = new Error('Project name already exists');
        error.status = 409;
        throw error;
      }
    }

    if (updates.apiKey) {
      updates.apiKey = encrypt(updates.apiKey);
    }

    return await this.projectWriteRepository.update(id, updates);
  }

  async delete(id) {
    return await this.projectWriteRepository.delete(id);
  }
}
