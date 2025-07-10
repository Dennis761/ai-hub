class ProjectInputService {
    normalize(input) {
      return {
        name: input.name?.trim(),
        status: input.status?.trim(),
        apiKey: input.apiKey?.trim()
      };
    }
  
    validate({ name, status, apiKey }) {
      const allowedStatuses = ['active', 'inactive', 'archived'];
  
      if (!name || typeof name !== 'string' || name.length === 0) {
        throw new Error('Invalid or missing "name".');
      }
  
      if (status && !allowedStatuses.includes(status)) {
        throw new Error('Invalid "status". Allowed: active, inactive, archived.');
      }
  
      if (!apiKey || typeof apiKey !== 'string' || apiKey.length < 8) {
        throw new Error('Invalid or missing "apiKey". Must be at least 8 characters.');
      }
    }
  }
  
  export default ProjectInputService;  