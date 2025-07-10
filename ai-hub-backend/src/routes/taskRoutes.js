import express from 'express';
import checkAuth from '../middlewares/authMiddleware.js';
import taskContainer from '../di/task/rootTaskContainer.js';

const router = express.Router();

const { taskReadController, taskWriteController, taskStatusController } = taskContainer;

// Read
router.get('/project/:projectId', checkAuth, taskReadController.getTasksByProject.bind(taskReadController));
router.get('/:id', checkAuth, taskReadController.getTaskById.bind(taskReadController));

// Write
router.post('/', checkAuth, taskWriteController.create.bind(taskWriteController));
router.patch('/:id', checkAuth, taskWriteController.update.bind(taskWriteController));
router.delete('/:id', checkAuth, taskWriteController.delete.bind(taskWriteController));

// Status
router.post('/:id/status', checkAuth, taskStatusController.setStatus.bind(taskStatusController));

export default router;