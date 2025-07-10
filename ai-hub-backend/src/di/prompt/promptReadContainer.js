import PromptReadController from '../../controllers/prompt/promptReadController.js';

import PromptReadRepository from '../../repositories/prompt/promptReadRepository.js';
import PromptReadManager from '../../managers/prompt/promptReadManager.js';
import PromptReadService from '../../use-cases/prompt/promptReadService.js';
 
// Repositories
const promptReadRepository = new PromptReadRepository();

// Managers
const promptReadManager = new PromptReadManager(promptReadRepository);

// Services
const promptReadService = new PromptReadService(promptReadManager);

const promptReadController = new PromptReadController(promptReadService);

export default promptReadController;
