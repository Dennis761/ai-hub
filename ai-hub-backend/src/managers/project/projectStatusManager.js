export default class ProjectStatusManager {
    constructor(projectWriteRepository) {
      this.projectWriteRepository = projectWriteRepository;
    }
  
    async setStatus(id, status) {
      return await this.projectWriteRepository.update(id, { status });
    }
  }
  