class PromptOrderManager {
    constructor(promptOrderRepository) {
      this.promptOrderRepository = promptOrderRepository;
    }
  
    async reorder(orderArray) {
      return this.promptOrderRepository.bulkUpdateExecutionOrder(orderArray);
    }
  }
  
  export default PromptOrderManager;
  