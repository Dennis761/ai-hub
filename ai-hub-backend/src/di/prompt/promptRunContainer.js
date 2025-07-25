import PromptRunController from '../../controllers/prompt/promptRunController.js';

import PromptReadRepository from '../../repositories/prompt/promptReadRepository.js';
import PromptWriteRepository from '../../repositories/prompt/promptWriteRepository.js';
import PromptReadManager from '../../managers/prompt/promptReadManager.js';
import PromptHistoryManager from '../../managers/prompt/promptHistoryManager.js';
import PromptReadService from '../../use-cases/prompt/promptReadService.js';
import PromptHistoryService from '../../use-cases/prompt/promptHistoryService.js';
import PromptInputService from '../../use-cases/prompt/promptInputService.js';
import PromptRunService from '../../use-cases/prompt/promptRunService.js';

import TaskReadRepository from '../../repositories/task/taskReadRepository.js';
import TaskReadManager from '../../managers/task/taskReadManager.js';
import TaskReadService from '../../use-cases/task/taskReadService.js';

import ApiKeyReadRepository from '../../repositories/apiKey/apiKeyReadRepository.js';
import ApiKeyWriteRepository from '../../repositories/apiKey/apiKeyWriteRepository.js';
import ApiKeyCountRepository from '../../repositories/apiKey/apiKeyCountRepository.js';

import ApiKeyReadManager from '../../managers/apiKey/apiKeyReadManager.js';
import ApiKeyWriteManager from '../../managers/apiKey/apiKeyWriteManager.js';

import ApiKeyReadService from '../../use-cases/apiKey/apiKeyReadService.js';
import ApiKeyWriteService from '../../use-cases/apiKey/apiKeyWriteService.js';

import ProjectReadRepository from '../../repositories/project/projectReadRepository.js';
import ProjectReadManager from '../../managers/project/projectReadManager.js';
import ProjectReadService from '../../use-cases/project/projectReadService.js';

// -------------------------
// Prompt layer setup
// -------------------------

const promptReadRepository = new PromptReadRepository();
const promptWriteRepository = new PromptWriteRepository();

const promptReadManager = new PromptReadManager(promptReadRepository);
const promptHistoryManager = new PromptHistoryManager({
  promptReadRepository,
  promptWriteRepository,
});

const promptReadService = new PromptReadService(promptReadManager);
const promptInputService = new PromptInputService();

// -------------------------
// Task layer setup
// -------------------------

const taskReadRepository = new TaskReadRepository();
const taskReadManager = new TaskReadManager(taskReadRepository);
const taskReadService = new TaskReadService(taskReadManager);

// -------------------------
// Project layer setup
// -------------------------

const projectReadRepository = new ProjectReadRepository();
const projectReadManager = new ProjectReadManager(projectReadRepository);
const projectReadService = new ProjectReadService(projectReadManager);

// -------------------------
// ApiKey layer setup
// -------------------------

const apiKeyReadRepository = new ApiKeyReadRepository();
const apiKeyWriteRepository = new ApiKeyWriteRepository();
const apiKeyCountRepository = new ApiKeyCountRepository();

const apiKeyReadManager = new ApiKeyReadManager(apiKeyReadRepository);
const apiKeyWriteManager = new ApiKeyWriteManager({
  apiKeyReadRepository,
  apiKeyWriteRepository,
  apiKeyCountRepository,
});

const apiKeyReadService = new ApiKeyReadService(apiKeyReadManager);
const apiKeyWriteService = new ApiKeyWriteService(apiKeyWriteManager);

// -------------------------
// PromptHistoryService setup
// -------------------------

const promptHistoryService = new PromptHistoryService({
  promptHistoryManager,
  promptReadManager,
  taskReadManager,
  projectReadService
});

const promptRunService = new PromptRunService({
  promptReadService,
  taskReadService,
  apiKeyWriteService,
  apiKeyReadService,
  promptInputService,
  promptHistoryService,
});

const promptRunController = new PromptRunController(promptRunService);

export default promptRunController;
