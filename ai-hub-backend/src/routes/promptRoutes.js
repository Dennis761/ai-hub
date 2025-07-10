import express from 'express';
import promptContainer from '../di/prompt/rootPromptContainer.js';
import checkAuth from '../middlewares/authMiddleware.js';

const router = express.Router();

const {
  promptWriteController,
  promptRunController,
  promptReadController,
  promptOrderController,
  promptHistoryController
} = promptContainer;

// Write
router.post('/', checkAuth, promptWriteController.create.bind(promptWriteController));
router.patch('/:id', checkAuth, promptWriteController.update.bind(promptWriteController));
router.delete('/:id', checkAuth, promptWriteController.delete.bind(promptWriteController));

// Run
router.post('/:id/run', checkAuth, promptRunController.run.bind(promptRunController));

// Read
router.get('/task/:taskId', checkAuth, promptReadController.getPromptsByTask.bind(promptReadController));

// Order
router.post('/reorder', checkAuth, promptOrderController.reorder.bind(promptOrderController));

// History
router.post('/:id/rollback', checkAuth, promptHistoryController.rollback.bind(promptHistoryController));

export default router;
