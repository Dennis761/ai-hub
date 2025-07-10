import ApiKeyReadRepository from '../../repositories/apiKey/apiKeyReadRepository.js';
import ApiKeyReadManager from '../../managers/apiKey/apiKeyReadManager.js';
import ApiKeyReadService from '../../use-cases/apiKey/apiKeyReadService.js';
import ApiKeyReadController from '../../controllers/apiKey/apiKeyReadController.js';

// Repositories
const apiKeyReadRepository = new ApiKeyReadRepository();

// Managers
const apiKeyReadManager = new ApiKeyReadManager(apiKeyReadRepository);

// Services
const apiKeyReadService = new ApiKeyReadService(apiKeyReadManager);

const apiKeyReadController = new ApiKeyReadController(apiKeyReadService);

export default apiKeyReadController