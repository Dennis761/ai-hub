import ProjectReadRepository from '../../repositories/project/projectReadRepository.js';
import ProjectWriteRepository from '../../repositories/project/projectWriteRepository.js';

import ProjectAccessManager from '../../managers/project/projectAccessManager.js';
import ProjectAccessService from '../../use-cases/project/projectAccessService.js';
import ProjectReadController from '../../controllers/project/projectReadController.js';

// Repositories
const projectReadRepository = new ProjectReadRepository();
const projectWriteRepository = new ProjectWriteRepository();

// Managers
const projectAccessManager = new ProjectAccessManager(projectReadRepository, projectWriteRepository);

// Services
const projectAccessService = new ProjectAccessService(projectAccessManager);

const projectReadController = new ProjectReadController(projectAccessService);

export default projectReadController