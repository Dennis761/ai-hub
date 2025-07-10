import ProjectReadRepository from '../../repositories/project/projectReadRepository.js';
import ProjectWriteRepository from '../../repositories/project/projectWriteRepository.js';

import ProjectWriteManager from '../../managers/project/projectWriteManager.js';
import ProjectWriteService from '../../use-cases/project/projectWriteService.js';
import ProjectInputService from '../../use-cases/project/projectInputService.js';
import ProjectWriteController from '../../controllers/project/projectWriteController.js';

// Reposities
const projectReadRepository = new ProjectReadRepository();
const projectWriteRepository = new ProjectWriteRepository();

// Managers
const projectWriteManager = new ProjectWriteManager({ projectWriteRepository, projectReadRepository });

// Services
const projectWriteService = new ProjectWriteService(projectWriteManager);
const projectInputService = new ProjectInputService();

const projectWriteController = new ProjectWriteController({
  projectWriteService,
  projectInputService,
});

export default projectWriteController