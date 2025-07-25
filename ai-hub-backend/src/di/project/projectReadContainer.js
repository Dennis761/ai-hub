import ProjectReadRepository from '../../repositories/project/projectReadRepository.js';

import ProjectReadManager from '../../managers/project/projectReadManager.js';
import ProjectReadService from '../../use-cases/project/projectReadService.js';
import ProjectReadController from '../../controllers/project/projectReadController.js';

// Repositories
const projectReadRepository = new ProjectReadRepository();

// Managers
const projectReadManager = new ProjectReadManager(projectReadRepository);

// Services
const projectReadService = new ProjectReadService(projectReadManager);
 
const projectReadController = new ProjectReadController(projectReadService);

export default projectReadController