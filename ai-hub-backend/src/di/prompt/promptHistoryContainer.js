import PromptHistoryController from '../../controllers/prompt/promptHistoryController.js';

import PromptReadRepository from '../../repositories/prompt/promptReadRepository.js';
import PromptWriteRepository from '../../repositories/prompt/promptWriteRepository.js';

import PromptHistoryManager from '../../managers/prompt/promptHistoryManager.js';
import PromptHistoryService from '../../use-cases/prompt/promptHistoryService.js';

// Repositories
const promptReadRepository = new PromptReadRepository();
const promptWriteRepository = new PromptWriteRepository();

// Managers
const promptHistoryManager = new PromptHistoryManager({ promptReadRepository, promptWriteRepository });

// Services
const promptHistoryService = new PromptHistoryService(promptHistoryManager);

const promptHistoryController = new PromptHistoryController(promptHistoryService);

export default promptHistoryController
