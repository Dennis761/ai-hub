import {
  incrementProjectEditCount,
  isInTop100Projects,
  cacheProject
} from '../../utils/cache/projectCache.js';

class TaskWriteService {
  constructor({ taskWriteManager, projectReadService }) {
    this.taskWriteManager = taskWriteManager;
    this.projectReadService = projectReadService;
  }
 
  async create(data) {
    const task = await this.taskWriteManager.create(data);

    if (task.projectId) {
      // Track edit activity by incrementing the project's edit count
      await incrementProjectEditCount(task.projectId);

      // If the project is among the top 100, update its cache
      if (await isInTop100Projects(task.projectId)) { 
        const project = await this.projectReadService.getById(task.projectId);

        // Cache the updated project data
        await cacheProject(task.projectId, project);
      } 
    }

    return task;
  }

  async update(id, updates) {
    const updated = await this.taskWriteManager.update(id, updates);
 
    if (updated?.projectId) {
      // Track edit activity by incrementing the project's edit count
      await incrementProjectEditCount(updated.projectId);

      // If the project is among the top 100, update its cache
      if (await isInTop100Projects(updated.projectId)) {
        const project = await this.projectReadService.getById(updated.projectId);

        // Cache the updated project data
        await cacheProject(updated.projectId, project);
      }
    }

    return updated;
  }

  async delete(id) {
    const deleted = await this.taskWriteManager.delete(id);

    if (deleted?.projectId) {
      // Track edit activity by incrementing the project's edit count
      await incrementProjectEditCount(deleted.projectId);

      // If the project is among the top 100, update its cache
      if (await isInTop100Projects(deleted.projectId)) {
        const project = await this.projectReadService.getById(deleted.projectId);

        // Cache the updated project data
        await cacheProject(deleted.projectId, project);
      }
    }

    return deleted;
  }
}

export default TaskWriteService;