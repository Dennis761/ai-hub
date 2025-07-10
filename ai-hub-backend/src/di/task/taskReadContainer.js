import TaskReadRepository from '../../repositories/task/taskReadRepository.js';
import TaskReadManager from '../../managers/task/taskReadManager.js';
import TaskReadService from '../../use-cases/task/taskReadService.js';
import TaskReadController from '../../controllers/task/taskReadController.js';

// Repositories
const taskReadRepository = new TaskReadRepository();

// Managers
const taskReadManager = new TaskReadManager(taskReadRepository);

// Services
const taskReadService = new TaskReadService(taskReadManager);

const taskReadController = new TaskReadController(taskReadService);

export default taskReadController