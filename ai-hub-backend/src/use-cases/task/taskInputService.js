class TaskInputService {
    normalize(input) {
      return {
        name: input.name?.trim(),
        description: input.description?.trim(),
        projectId: input.projectId?.trim(),
        apiMethod: input.apiMethod?.trim(),
        status: input.status?.trim()
      };
    }
  
    validate({ name, projectId, apiMethod }) {
      if (!name || typeof name !== 'string' || name.length === 0) {
        throw new Error('Invalid or missing "name".');
      }
  
      if (!projectId || typeof projectId !== 'string') {
        throw new Error('Invalid or missing "projectId".');
      }
  
      if (!apiMethod || typeof apiMethod !== 'string' || !apiMethod.includes('?')) {
        throw new Error('Invalid or missing "apiMethod". It must include query parameters.');
      }
    }
  
    validateStatus(status) {
      const allowed = ['archived', 'active', 'inactive'];
      if (!allowed.includes(status)) {
        throw new Error(`Invalid status: must be one of [${allowed.join(', ')}]`);
      }
    }
  }
  
  export default TaskInputService;
  