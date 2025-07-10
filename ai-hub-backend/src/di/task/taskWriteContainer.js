import TaskWriteRepository from '../../repositories/task/taskWriteRepository.js';
import TaskWriteManager from '../../managers/task/taskWriteManager.js';
import TaskWriteService from '../../use-cases/task/taskWriteService.js';
import TaskInputService from '../../use-cases/task/taskInputService.js';
import TaskWriteController from '../../controllers/task/taskWriteController.js';

// Repositories
const taskWriteRepository = new TaskWriteRepository();

// Managers
const taskWriteManager = new TaskWriteManager(taskWriteRepository);

// Services
const taskWriteService = new TaskWriteService(taskWriteManager);
const taskInputService = new TaskInputService();

const taskWriteController = new TaskWriteController({taskWriteService, taskInputService});

export default taskWriteController