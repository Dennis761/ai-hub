import ApiKeyWriteRepository from '../../repositories/apiKey/apiKeyWriteRepository.js';
import ApiKeyReadRepository from '../../repositories/apiKey/apiKeyReadRepository.js';
import ApiKeyCountRepository from '../../repositories/apiKey/apiKeyCountRepository.js';

import ApiKeyWriteManager from '../../managers/apiKey/apiKeyWriteManager.js';
import ApiKeyStatusManager from '../../managers/apiKey/apiKeyStatusManager.js';
import ApiKeyStatusService from '../../use-cases/apiKey/apiKeyStatusService.js';
import ApiKeyStatusController from '../../controllers/apiKey/apiKeyStatusController.js';

// Repositories
const apiKeyWriteRepository = new ApiKeyWriteRepository();
const apiKeyReadRepository = new ApiKeyReadRepository();
const apiKeyCountRepository = new ApiKeyCountRepository();

// Managers
const apiKeyWriteManager = new ApiKeyWriteManager({
  apiKeyReadRepository,
  apiKeyWriteRepository,
  apiKeyCountRepository,
});
const apiKeyStatusManager = new ApiKeyStatusManager(apiKeyWriteManager);

// Services
const apiKeyStatusService = new ApiKeyStatusService(apiKeyStatusManager);

const apiKeyStatusController = new ApiKeyStatusController(apiKeyStatusService);

export default apiKeyStatusController