import {
  incrementProjectEditCount,
  isInTop100Projects,
  cacheProject
} from '../../utils/cache/projectCache.js';

class PromptHistoryService {
  constructor({ promptHistoryManager, promptReadManager, taskReadManager, projectReadService }) {
    this.promptHistoryManager = promptHistoryManager;
    this.promptReadManager = promptReadManager;
    this.taskReadManager = taskReadManager;
    this.projectReadService = projectReadService;
  }

  async addToHistory(promptId, promptText, responseText) {
    return this.promptHistoryManager.addToHistory(promptId, promptText, responseText);
  }

  async rollbackToHistoryVersion(promptId, versionIndex) {
    const result = await this.promptHistoryManager.rollback(promptId, versionIndex);

    // Retrieve related task and project for caching updates
    const prompt = await this.promptReadManager.getById(promptId);
    const task = await this.taskReadManager.getById(prompt.taskId);

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

    return result;
  }
}

export default PromptHistoryService;
