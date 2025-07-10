import ProjectWriteRepository from '../../repositories/project/projectWriteRepository.js';

import ProjectStatusManager from '../../managers/project/projectStatusManager.js';
import ProjectStatusService from '../../use-cases/project/projectStatusService.js';
import ProjectStatusController from '../../controllers/project/projectStatusController.js';

// Repositories
const projectWriteRepository = new ProjectWriteRepository();

//  Managers
const projectStatusManager = new ProjectStatusManager(projectWriteRepository);

// Services
const projectStatusService = new ProjectStatusService(projectStatusManager);

const projectStatusController = new ProjectStatusController(projectStatusService);

export default projectStatusController