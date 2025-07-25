import {
  incrementProjectEditCount,
  isInTop100Projects,
  cacheProject,
} from '../../utils/cache/projectCache.js';

export default class ProjectWriteService {
    constructor({projectWriteManager}) {
      this.projectWriteManager = projectWriteManager;
    }
  
    create(data) {
      return this.projectWriteManager.create(data);
    }
  
    async update(id, updates) {
      const updatedProject = await this.projectWriteManager.update(id, updates);

      // Track edit activity by incrementing the project's edit count
      await incrementProjectEditCount(id);

      // Cache handling: if the project is in the top 100, update the cache
      if (await isInTop100Projects(id)) {
        await cacheProject(id, updatedProject);
      }
    
      return updatedProject;
    }    
  
    delete(id) {
      return this.projectWriteManager.delete(id);
    }
  }
  