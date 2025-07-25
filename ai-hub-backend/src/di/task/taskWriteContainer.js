import TaskWriteRepository from '../../repositories/task/taskWriteRepository.js';
import TaskReadRepository from '../../repositories/task/taskReadRepository.js';
import TaskWriteManager from '../../managers/task/taskWriteManager.js';
import TaskWriteService from '../../use-cases/task/taskWriteService.js';
import TaskInputService from '../../use-cases/task/taskInputService.js';
import TaskWriteController from '../../controllers/task/taskWriteController.js';

import ProjectReadRepository from '../../repositories/project/projectReadRepository.js';
import ProjectReadManager from '../../managers/project/projectReadManager.js';
import ProjectReadService from '../../use-cases/project/projectReadService.js';

// Repositories
const taskWriteRepository = new TaskWriteRepository();
const taskReadRepository = new TaskReadRepository();
const projectReadRepository = new ProjectReadRepository();

// Managers
const projectReadManager = new ProjectReadManager(projectReadRepository);
const taskWriteManager = new TaskWriteManager({taskWriteRepository, taskReadRepository});

// Services
const projectReadService = new ProjectReadService(projectReadManager);
const taskWriteService = new TaskWriteService({ taskWriteManager, projectReadService});
const taskInputService = new TaskInputService();

const taskWriteController = new TaskWriteController({taskWriteService, taskInputService});

export default taskWriteController