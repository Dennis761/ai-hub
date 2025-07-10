import express from 'express';
import checkAuth from '../middlewares/authMiddleware.js';
import projectContainer from '../di/project/rootProjectContainer.js';

const router = express.Router();

const {
  projectAccessController,
  projectReadController,
  projectWriteController,
  projectStatusController
} = projectContainer;

// Read
router.get('/my', checkAuth, projectReadController.getMyProjects.bind(projectReadController));

// Write
router.post('/', checkAuth, projectWriteController.create.bind(projectWriteController));
router.patch('/:id', checkAuth, projectWriteController.update.bind(projectWriteController));
router.delete('/:id', checkAuth, projectWriteController.delete.bind(projectWriteController));

// Access
router.post('/join', checkAuth, projectAccessController.joinProject.bind(projectAccessController));

// Status
router.post('/:id/status', checkAuth, projectStatusController.setStatus.bind(projectStatusController));

export default router;
