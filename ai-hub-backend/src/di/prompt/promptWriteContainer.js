import PromptInputService from '../../use-cases/prompt/promptInputService.js';

import PromptWriteRepository from '../../repositories/prompt/promptWriteRepository.js';
import PromptWriteManager from '../../managers/prompt/promptWriteManager.js';
import PromptWriteService from '../../use-cases/prompt/promptWriteService.js';

import PromptHistoryManager from '../../managers/prompt/promptHistoryManager.js';
import PromptHistoryService from '../../use-cases/prompt/promptHistoryService.js';

import PromptReadRepository from '../../repositories/prompt/promptReadRepository.js';
import PromptReadManager from '../../managers/prompt/promptReadManager.js';
import PromptReadService from '../../use-cases/prompt/promptReadService.js';

import PromptWriteController from '../../controllers/prompt/promptWriteController.js';

import ApiKeyReadRepository from '../../repositories/apiKey/apiKeyReadRepository.js';
import ApiKeyReadManager from '../../managers/apiKey/apiKeyReadManager.js';
import ApiKeyReadService from '../../use-cases/apiKey/apiKeyReadService.js';

import TaskReadRepository from '../../repositories/task/taskReadRepository.js';
import TaskReadManager from '../../managers/task/taskReadManager.js';
import TaskReadService from '../../use-cases/task/taskReadService.js';

import ProjectReadRepository from '../../repositories/project/projectReadRepository.js';
import ProjectReadManager from '../../managers/project/projectReadManager.js';
import ProjectReadService from '../../use-cases/project/projectReadService.js';

// -------------------------
// Project layer setup
// -------------------------
const projectReadRepository = new ProjectReadRepository();
const projectReadManager = new ProjectReadManager(projectReadRepository);
const projectReadService = new ProjectReadService(projectReadManager);

// -------------------------
// Task layer setup
// -------------------------
const taskReadRepository = new TaskReadRepository();
const taskReadManager = new TaskReadManager(taskReadRepository);
const taskReadService = new TaskReadService(taskReadManager);

// -------------------------
// API Key layer setup
// -------------------------
const apiKeyReadRepository = new ApiKeyReadRepository();
const apiKeyReadManager = new ApiKeyReadManager(apiKeyReadRepository);
const apiKeyReadService = new ApiKeyReadService(apiKeyReadManager);

// -------------------------
// Prompt layer setup
// -------------------------
const promptWriteRepository = new PromptWriteRepository();
const promptReadRepository = new PromptReadRepository();

const promptWriteManager = new PromptWriteManager({
  promptWriteRepository,
  promptReadRepository,
});
const promptReadManager = new PromptReadManager(promptReadRepository);
const promptHistoryManager = new PromptHistoryManager({
  promptWriteRepository,
  promptReadRepository,
});

const promptWriteService = new PromptWriteService({
  promptWriteManager,
  taskReadManager,
  apiKeyReadService,
  projectReadService,
});
const promptReadService = new PromptReadService(promptReadManager);
const promptHistoryService = new PromptHistoryService({
  promptHistoryManager,
  promptReadManager,
  taskReadManager,
  projectReadService
});
const promptInputService = new PromptInputService();
 
const promptWriteController = new PromptWriteController({
  promptWriteService,
  promptReadService,
  promptHistoryService,
  promptInputService,
  taskReadService,
});

export default promptWriteController;