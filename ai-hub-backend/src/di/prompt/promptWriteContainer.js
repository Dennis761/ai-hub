import PromptWriteRepository from '../../repositories/prompt/promptWriteRepository.js';
import PromptReadRepository from '../../repositories/prompt/promptReadRepository.js';
import ApiKeyReadRepository from '../../repositories/apiKey/apiKeyReadRepository.js';
import TaskReadRepository from '../../repositories/task/taskReadRepository.js';

import PromptWriteManager from '../../managers/prompt/promptWriteManager.js';
import PromptReadManager from '../../managers/prompt/promptReadManager.js';
import PromptHistoryManager from '../../managers/prompt/promptHistoryManager.js';
import ApiKeyReadManager from '../../managers/apiKey/apiKeyReadManager.js';
import TaskReadManager from '../../managers/task/taskReadManager.js';

import PromptWriteService from '../../use-cases/prompt/promptWriteService.js';
import PromptReadService from '../../use-cases/prompt/promptReadService.js';
import PromptHistoryService from '../../use-cases/prompt/promptHistoryService.js';
import PromptInputService from '../../use-cases/prompt/promptInputService.js';
import TaskReadService from '../../use-cases/task/taskReadService.js';
import ApiKeyReadService from '../../use-cases/apiKey/apiKeyReadService.js';

import PromptWriteController from '../../controllers/prompt/promptWriteController.js';

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

const promptWriteManager = new PromptWriteManager({ promptWriteRepository, promptReadRepository });
const promptReadManager = new PromptReadManager(promptReadRepository);
const promptHistoryManager = new PromptHistoryManager({ promptWriteRepository, promptReadRepository });

const promptWriteService = new PromptWriteService({ promptWriteManager, apiKeyReadService });
const promptReadService = new PromptReadService(promptReadManager);
const promptHistoryService = new PromptHistoryService(promptHistoryManager);
const promptInputService = new PromptInputService();


// -------------------------
// Task layer setup
// -------------------------
const taskReadRepository = new TaskReadRepository();
const taskReadManager = new TaskReadManager(taskReadRepository);
const taskReadService = new TaskReadService(taskReadManager);

const promptWriteController = new PromptWriteController({
  promptWriteService,
  promptReadService,
  promptHistoryService,
  promptInputService,
  taskReadService,
});

export default promptWriteController