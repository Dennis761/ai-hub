import ApiKeyReadRepository from '../../repositories/apiKey/apiKeyReadRepository.js';
import ApiKeyWriteRepository from '../../repositories/apiKey/apiKeyWriteRepository.js';
import ApiKeyCountRepository from '../../repositories/apiKey/apiKeyCountRepository.js';

import ApiKeyWriteManager from '../../managers/apiKey/apiKeyWriteManager.js';
import ApiKeyWriteService from '../../use-cases/apiKey/apiKeyWriteService.js';
import ApiKeyInputService from '../../use-cases/apiKey/apiKeyInputService.js';
import ApiKeyWriteController from '../../controllers/apiKey/apiKeyWriteController.js';

// Repositories
const apiKeyReadRepository = new ApiKeyReadRepository();
const apiKeyWriteRepository = new ApiKeyWriteRepository();
const apiKeyCountRepository = new ApiKeyCountRepository();

// Managers
const apiKeyWriteManager = new ApiKeyWriteManager({
  apiKeyReadRepository,
  apiKeyWriteRepository,
  apiKeyCountRepository,
});

// Services
const apiKeyWriteService = new ApiKeyWriteService(apiKeyWriteManager);
const apiKeyInputService = new ApiKeyInputService();

const apiKeyWriteController = new ApiKeyWriteController({
  apiKeyWriteService,
  apiKeyInputService,
}); 

export default apiKeyWriteController