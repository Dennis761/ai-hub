import ProjectReadRepository from '../../repositories/project/projectReadRepository.js';

import ProjectAccessManager from '../../managers/project/projectAccessManager.js';
import ProjectAccessService from '../../use-cases/project/projectAccessService.js';
import ProjectAccessController from '../../controllers/project/projectAccessController.js';

// Repositories
const projectReadRepository = new ProjectReadRepository();

// Managers
const projectAccessManager = new ProjectAccessManager(projectReadRepository);

// Services
const projectAccessService = new ProjectAccessService(projectAccessManager);

const projectAccessController = new ProjectAccessController(projectAccessService);

export default projectAccessController