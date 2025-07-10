class PromptOrderService {
  constructor(promptOrderManager) {
    this.promptOrderManager = promptOrderManager;
  }

  async reorder(orderList) {
    return this.promptOrderManager.reorder(orderList);
  }
}

export default PromptOrderService;
