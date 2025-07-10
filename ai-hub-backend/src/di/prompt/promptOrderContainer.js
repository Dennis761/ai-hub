import PromptOrderController from '../../controllers/prompt/promptOrderController.js';

import PromptOrderRepository from '../../repositories/prompt/promptOrderRepository.js';

import PromptOrderManager from '../../managers/prompt/promptOrderManager.js';
import PromptOrderService from '../../use-cases/prompt/promptOrderService.js';

// Repositories
const promptOrderRepository = new PromptOrderRepository();

// Managers
const promptOrderManager = new PromptOrderManager(promptOrderRepository);

// Services
const promptOrderService = new PromptOrderService(promptOrderManager);

const promptOrderController = new PromptOrderController(promptOrderService);

export default promptOrderController
