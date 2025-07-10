class PromptOrderController {
    constructor( promptOrderService ) {
      this.promptOrderService = promptOrderService;
    }
  
    async reorder(req, res) {
      try {
        const result = await this.promptOrderService.reorder(req.body);
        res.json({ message: 'Order updated', result });
      } catch (error) {
        res.status(500).json({ error: 'Failed to reorder prompts', details: error.message });
      }
    }
  }
  
  export default PromptOrderController;
  