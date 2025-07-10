import ProjectReadRepository from '../../repositories/project/projectReadRepository.js';
import ProjectWriteRepository from '../../repositories/project/projectWriteRepository.js';

import ProjectAccessManager from '../../managers/project/projectAccessManager.js';
import ProjectAccessService from '../../use-cases/project/projectAccessService.js';
import ProjectAccessController from '../../controllers/project/projectAccessController.js';

// Repositories
const projectReadRepository = new ProjectReadRepository();
const projectWriteRepository = new ProjectWriteRepository();

// Managers
const projectAccessManager = new ProjectAccessManager(projectReadRepository, projectWriteRepository);

// Services
const projectAccessService = new ProjectAccessService(projectAccessManager);

const projectAccessController = new ProjectAccessController(projectAccessService);

export default projectAccessController