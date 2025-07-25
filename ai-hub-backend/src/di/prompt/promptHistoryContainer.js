import ProjectReadRepository from '../../repositories/project/projectReadRepository.js';
import ProjectReadManager from '../../managers/project/projectReadManager.js';
import ProjectReadService from '../../use-cases/project/projectReadService.js';

import TaskReadRepository from '../../repositories/task/taskReadRepository.js';
import TaskReadManager from '../../managers/task/taskReadManager.js';

import PromptReadRepository from '../../repositories/prompt/promptReadRepository.js';
import PromptReadManager from '../../managers/prompt/promptReadManager.js';
import PromptWriteRepository from '../../repositories/prompt/promptWriteRepository.js';
import PromptHistoryManager from '../../managers/prompt/promptHistoryManager.js';
import PromptHistoryService from '../../use-cases/prompt/promptHistoryService.js';

import PromptHistoryController from '../../controllers/prompt/promptHistoryController.js';

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


// -------------------------
// Prompt layer setup
// -------------------------
const promptReadRepository = new PromptReadRepository();
const promptReadManager = new PromptReadManager(promptReadRepository);
const promptWriteRepository = new PromptWriteRepository();
const promptHistoryManager = new PromptHistoryManager({ promptReadRepository, promptWriteRepository });
const promptHistoryService = new PromptHistoryService({
    promptHistoryManager,
    promptReadManager,
    taskReadManager,
    projectReadService
  });

const promptHistoryController = new PromptHistoryController(promptHistoryService);

export default promptHistoryController
