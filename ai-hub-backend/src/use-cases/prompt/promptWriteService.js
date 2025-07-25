import {
  incrementProjectEditCount,
  isInTop100Projects,
  cacheProject
} from '../../utils/cache/projectCache.js';

class PromptWriteService {
  constructor({ promptWriteManager, taskReadManager, apiKeyReadService, projectReadService }) {
    this.promptWriteManager = promptWriteManager;
    this.taskReadManager = taskReadManager;
    this.apiKeyReadService = apiKeyReadService;
    this.projectReadService = projectReadService;
  }

  async create(data) {
    const prompt = await this.promptWriteManager.create(data);

    // Retrieve related task and project for caching updates
    const task = await this.taskReadManager.getById(prompt.taskId);

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

    const { modelName, modelProvider } = await this.apiKeyReadService.getModelById(prompt.modelId);

    return { prompt, modelName, modelProvider };
  }

  async update(promptId, updates) {
    const updatedPrompt = await this.promptWriteManager.update(promptId, updates);

    // Retrieve related task and project for caching updates
    const task = await this.taskReadManager.getById(updatedPrompt.taskId);

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

    return updatedPrompt;
  }

  async delete(promptId) {
    const deletedPrompt = await this.promptWriteManager.delete(promptId);

    // Retrieve related task and project for caching updates
    const task = await this.taskReadManager.getById(deletedPrompt.taskId);
    
    if (task?.projectId) {
      // Track edit activity by incrementing the project's edit count
      await incrementProjectEditCount(task.projectId);

      // If the project is among the top 100, update its cache
      if (await isInTop100Projects(task.projectId)) {
        const project = await this.projectReadService.getById(task.projectId);

        // Cache the updated project data
        await cacheProject(task.projectId, project);
      }
    }

    return deletedPrompt;
  }
}

export default PromptWriteService;
