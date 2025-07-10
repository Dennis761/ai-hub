import TaskWriteRepository from '../../repositories/task/taskWriteRepository.js';
import TaskStatusManager from '../../managers/task/taskStatusManager.js';
import TaskStatusService from '../../use-cases/task/taskStatusService.js';
import TaskInputService from '../../use-cases/task/taskInputService.js';
import TaskStatusController from '../../controllers/task/taskStatusController.js';

// Repositories
const taskWriteRepository = new TaskWriteRepository();

// Managers
const taskStatusManager = new TaskStatusManager(taskWriteRepository);

// Services
const taskInputService = new TaskInputService()
const taskStatusService = new TaskStatusService(taskStatusManager);

const taskStatusController = new TaskStatusController({taskStatusService, taskInputService});

export default taskStatusController